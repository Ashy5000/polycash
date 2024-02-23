package main

import (
	"crypto/dsa"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

var lostBlock = false

func CreateBlock(sender dsa.PublicKey, recipient dsa.PublicKey, amount float64) (Block, error) {
	start := time.Now()
	previousBlock, previousBlockFound := GetLastMinedBlock()
	if !previousBlockFound {
		previousBlock.Difficulty = 100000
		previousBlock.MiningTime = time.Minute
	}
	block := Block{
		Miner:      GetKey().PublicKey,
		Sender:     sender,
		Recipient:  recipient,
		Amount:     amount,
		Nonce:      0,
		Difficulty: previousBlock.Difficulty * (60 / uint64(previousBlock.MiningTime.Seconds())),
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	fmt.Printf("Mining block with difficulty %d\n", block.Difficulty)
	for hash > 9223372036854776000/block.Difficulty {
		if lostBlock {
			lostBlock = false
			return Block{}, errors.New("lost block")
		} else {
			block.Nonce++
			hashBytes = HashBlock(block)
			hash = binary.BigEndian.Uint64(hashBytes[:])
		}
	}
	block.MiningTime = time.Since(start)
	return block, nil
}
