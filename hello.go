package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}

//aes的加密字符串
var aesKey = "astaxie12798akljzmknm.ahkjkljlaa"

func main() {
	fmt.Println("Hello World")

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/login", func(c *gin.Context) {
			session, err := EncryptSecureData([]byte("I'm Zuolar"))
			if err != nil {
				c.JSON(http.StatusOK, "Login Failed ---> "+err.Error())
				return
			}
			cookie := &http.Cookie{
				Name:  "session",
				Value: session,
			}
			http.SetCookie(c.Writer, cookie)
			c.JSON(http.StatusOK, "Login Done")
		})

		api.GET("user", func(c *gin.Context) {
			session, err := c.Cookie("session")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, "login first")
				return
			}

			data, err := DecyptSecureData(session)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusOK, "login first")
				return
			}

			if string(data) != "I'm Zuolar" {
				c.AbortWithStatusJSON(http.StatusOK, "login first")
				return
			}

			c.Set("name", time.Now().Unix())
			a := c.Query("a")
			fmt.Println("Get Query a --->", a)
			c.Next()
			cd := c.Query("c")
			fmt.Println("Get Query c --->", cd)
		}, func(c *gin.Context) {
			name := c.GetInt64("name")
			b := c.Query("b")
			fmt.Println("Get Query b --->", b)
			c.JSON(http.StatusOK, name)
			return
		})
	}

	err := r.Run(":8000")
	log.Fatal(err)
}

// EncryptSecureData 加密敏感資料
func EncryptSecureData(data []byte) (string, error) {
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return "", err
	}

	//加密字符串
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	session := make([]byte, len(data))
	cfb.XORKeyStream(session, data)
	fmt.Printf("%s=>%x\n", string(data), session)
	return fmt.Sprintf("%x", session), nil
}

// DecyptSecureData 解密敏感資料
func DecyptSecureData(session string) ([]byte, error) {
	s, err := hex.DecodeString(session)
	// 创建加密算法aes
	c, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		return nil, err
	}

	// 解密字符串
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	data := make([]byte, len(s))
	cfbdec.XORKeyStream(data, s)
	fmt.Printf("%x=>%s\n", s, data)
	return data, nil
}
