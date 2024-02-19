package main

import (
	"crypto/sha256"
	"fmt"
)

func HashBlock(block Block) [32]byte {
	blockBytes := []byte(fmt.Sprintf("%v", block))
	sum := sha256.Sum256(blockBytes)
	return sum
}
