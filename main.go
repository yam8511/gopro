package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var chain *BlockChain

func main() {
	fmt.Println("Hello World")

	chain = NewBlockChain()
	fmt.Println("鏈長度", chain.Len())
	chain.NewBlock("筆電", 31.12)
	chain.NewBlock("滑鼠", 10.21)
	fmt.Println("鏈長度", chain.Len())

	r := gin.Default()
	r.GET("/", getBlockHandler)
	r.POST("/", writeBlockHandler)
	r.PUT("/", putBlockHandler)
	r.Run(":8000")
}

func getBlockHandler(c *gin.Context) {
	for _, block := range chain.Chain {
		fmt.Println(block)
	}
	c.JSON(http.StatusOK, chain.Chain)
}

func writeBlockHandler(c *gin.Context) {
	var newBlock Block
	err := c.ShouldBindJSON(&newBlock)
	if err != nil {
		c.JSON(http.StatusOK, "參數錯誤")
		return
	}

	newBlock, err = chain.NewBlock(newBlock.Name, newBlock.Price)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"block": newBlock,
		"len":   chain.Len(),
	})
}

func putBlockHandler(c *gin.Context) {
	var newBlock Block
	err := c.ShouldBindJSON(&newBlock)
	if err != nil {
		c.JSON(http.StatusOK, "參數錯誤")
		return
	}
	newBlock, err = chain.NewBlock(newBlock.Name, newBlock.Price)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"block": newBlock,
		"len":   chain.Len(),
	})
}
