package main

import (
	"crypto/dsa"
	"encoding/binary"
	"errors"
)

var lostBlock = false

func CreateBlock(sender dsa.PublicKey, recipient dsa.PublicKey, amount float64) (Block, error) {
	block := Block{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Nonce:     0,
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:])
	var maxHash uint64
	maxHash = 0x1000000000000000
	for hash > maxHash {
		if lostBlock {
			lostBlock = false
			return Block{
				Sender:    dsa.PublicKey{},
				Recipient: dsa.PublicKey{},
				Miner:     dsa.PublicKey{},
				Amount:    0,
				Nonce:     0,
			}, errors.New("lost block")
		} else {
			block.Nonce++
			hashBytes = HashBlock(block)
			hash = binary.BigEndian.Uint64(hashBytes[:])
		}
	}
	return block, nil
}
