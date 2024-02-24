package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func HashBlock(block Block) [32]byte {
	block.MiningTime = time.Minute
	blockBytes := []byte(fmt.Sprintf("%v", block))
	sum := sha256.Sum256(blockBytes)
	return sum
}
