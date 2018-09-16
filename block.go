package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// Block 區塊
type Block struct {
	// 區塊基本資訊
	Index        uint64
	Timestamp    string
	PreviousHash string
	Hash         string
	// 區塊資料內容
	Name  string
	Price float64
}

func (b Block) String() string {
	return fmt.Sprintf(`區塊 #%d :
	名稱: %s
	金額: %.2f
	時間: %s
	Hash: %s
	前一個 Hash: %s
`,
		b.Index,
		b.Name,
		b.Price,
		b.Timestamp,
		b.Hash,
		b.PreviousHash,
	)
}

// NewBlock 建立一個新的區塊
func NewBlock(index uint64, oldHash, name string, price float64) Block {
	newBlock := Block{
		Index:        index,
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		PreviousHash: oldHash,
		Name:         name,
		Price:        price,
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func calculateHash(block Block) string {
	hash := fmt.Sprintf("index=%d&ts=%s&prevHash=%s&name=%s&price=%f",
		block.Index,
		block.Timestamp,
		block.PreviousHash,
		block.Name,
		block.Price,
	)
	hash = fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))
	return hash
}
