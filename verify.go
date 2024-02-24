package main

import (
	"crypto/dsa"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
)

func VerifyTransaction(senderKey dsa.PublicKey, recipientKey dsa.PublicKey, amount string, r big.Int, s big.Int) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	isValid := dsa.Verify(&senderKey, []byte(fmt.Sprintf("%s:%s:%s", senderKey.Y, recipientKey.Y, strconv.FormatFloat(amountFloat, 'f', -1, 64))), &r, &s)
	return isValid
}

func VerifyBlock(block Block) bool {
	if !VerifyTransaction(block.Sender, block.Recipient, strconv.FormatFloat(block.Amount, 'f', -1, 64), block.R, block.S) {
		return false
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	if hash > 9223372036854776000/block.Difficulty {
		fmt.Println("Block has invalid hash. Ignoring block request.")
		fmt.Printf("Actual hash: %d\n", hash)
		return false
	}
	return true
}
