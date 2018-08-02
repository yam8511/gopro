package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Login 登入資料
type Login struct {
	Username string `json:"username" example:"zuolar"`
	Password string `json:"password" example:"qwe123"`
}

// User 使用者資料
type User struct {
	ID       int    `json:"id" example:"99"`
	Username string `json:"username" example:"zuolar"`
	Password string `json:"password" example:"qwe123"`
}

// @Summary 登入
// @Description 設定cookie
// @Tags user
// @Accept  json
// @Produce  json
// @Param  body  body  main.Login  true   "登入資訊"
// @Success 200 {string} string "登入狀況"
// @Router /api/login [post]
func login(c *gin.Context) {
	// 接受登入參數
	var input Login
	err := json.NewDecoder(c.Request.Body).Decode(&input)

	if err != nil {
		c.JSON(http.StatusOK, "登入失敗 ---> "+err.Error())
		return
	}

	if input.Password != "qwe123" {
		c.JSON(http.StatusOK, "登入失敗 ---> 密碼是 qwe123")
		return
	}

	cookie := &http.Cookie{
		Name:  "test_session",
		Value: input.Username + "@@@" + input.Password,
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, "登入成功")
}

// @Summary 變更密碼
// @Description 變更自己的密碼
// @Tags user
// @Accept  json
// @Produce  json
// @Param  password  query  string  true  "密碼資料"
// @Success 200 {string} string "變更狀態"
// @Failure 401 {string} string "請先登入"
// @Failure 500 {string} string "系統錯誤"
// @Router /api/password [put]
func updatePasswd(c *gin.Context) {
	newPasswd := c.Query("password")
	if newPasswd == "" {
		c.JSON(http.StatusOK, "變更失敗 ---> 參數 passord 必填")
		return
	}

	session, err := c.Cookie("test_session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "請先登入 ---> /api/login")
		return
	}

	data := strings.Split(session, "@@@")
	if len(data) != 2 {
		c.JSON(http.StatusInternalServerError, "系統錯誤 ---> 請稍後再嘗試")
		return
	}

	cookie := &http.Cookie{
		Name:  "test_session",
		Value: data[0] + "@@@" + newPasswd,
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, "變更成功 ---> "+newPasswd)
}

// @Summary 使用者資料
// @Description 取使用者ID
// @Tags user
// @Accept  json
// @Produce  json
// @Param  id  path  int  true  "使用者ID""
// @Success 200 {object} main.User "使用者資料"
// @Failure 401 {string} string "請先登入"
// @Router /api/user/{id} [get]
func getUser(c *gin.Context) {
	session, err := c.Cookie("test_session")
	if err != nil {
		c.JSON(http.StatusUnauthorized, "請先登入 ---> /api/login")
		return
	}

	data := strings.Split(session, "@@@")
	if len(data) < 2 {
		c.JSON(http.StatusUnauthorized, "請先登入 ---> /api/login")
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusOK, "取得失敗 ---> 無效ID")
		return
	}

	c.JSON(http.StatusOK, User{
		ID:       id,
		Username: data[0],
		Password: data[1],
	})
}

// @Summary 登出
// @Description 清除cookie
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {string} string "登出狀況"
// @Router /api/logout [delete]
func logout(c *gin.Context) {
	cookie := &http.Cookie{
		Name:    "test_session",
		Expires: time.Now(),
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, "登出成功")
}
