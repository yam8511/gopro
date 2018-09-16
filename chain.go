package main

import (
	"errors"
	"sync"
)

// NewBlockChain 建立一個新的區塊
func NewBlockChain() *BlockChain {
	newChain := new(BlockChain)
	return newChain
}

// BlockChain 區塊鏈
type BlockChain struct {
	Chain []Block
	Mx    sync.RWMutex
}

// Len 鏈的長度
func (c *BlockChain) Len() int {
	return len(c.Chain)
}

// NewBlock 新增區塊
func (c *BlockChain) NewBlock(name string, price float64) (Block, error) {
	c.Mx.Lock()
	defer c.Mx.Unlock()
	len := c.Len()
	if len == 0 {
		newBlock := NewBlock(0, "", name, price)
		c.Chain = append(c.Chain, newBlock)
		return newBlock, nil
	}

	previousBlock := c.LatestBlock()
	newBlock := NewBlock(previousBlock.Index+1, previousBlock.Hash, name, price)
	if !isBlockValid(newBlock, previousBlock) {
		return Block{}, errors.New("無效區塊")
	}
	c.Chain = append(c.Chain, newBlock)
	return newBlock, nil
}

// ReplaceChain 更換區塊鏈
func (c *BlockChain) ReplaceChain(newBlocks []Block) {
	if len(newBlocks) > len(c.Chain) {
		c.Chain = newBlocks
	}
}

// LatestBlock 最新的區塊
func (c *BlockChain) LatestBlock() Block {
	lastIndex := c.Len() - 1
	return c.Chain[lastIndex]
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreviousHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}
