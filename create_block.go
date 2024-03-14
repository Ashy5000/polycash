// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

// transactionHashes is a map of transaction hashes to their current status. 0 means the transaction is unmined, 1 means the transaction is being mined, and 2 means the transaction has been mined.
var transactionHashes = make(map[[32]byte]int)

func CreateBlock(sender dsa.PublicKey, recipient dsa.PublicKey, amount float64, r big.Int, s big.Int, transactionHash [32]byte, timestamp string) (Block, error) {
	start := time.Now()
	previousBlock, previousBlockFound := GetLastMinedBlock()
	if !previousBlockFound {
		previousBlock.Difficulty = 100000
		previousBlock.MiningTime = time.Minute
	}
	timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	block := Block{
		Miner:      GetKey().PublicKey,
		Sender:     sender,
		Recipient:  recipient,
		Amount:     amount,
		R:          r,
		S:          s,
		Nonce:      0,
		Difficulty: previousBlock.Difficulty * (60 / uint64(previousBlock.MiningTime.Seconds())),
		Timestamp:  time.Unix(0, timestampInt64),
	}
	if block.Difficulty < 10000 {
		block.Difficulty = 10000
	}
	if len(blockchain) > 0 {
		block.PreviousBlockHash = HashBlock(blockchain[len(blockchain)-1])
	} else {
		block.PreviousBlockHash = [32]byte{}
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	fmt.Printf("Mining block with difficulty %d\n", block.Difficulty)
	for hash > 9223372036854776000/block.Difficulty {
		if transactionHashes[transactionHash] == 2 {
			return Block{}, errors.New("lost block")
		} else {
			previousBlock, previousBlockFound = GetLastMinedBlock()
			if !previousBlockFound {
				previousBlock.Difficulty = 100000
				previousBlock.MiningTime = time.Minute
			}
			if len(blockchain) > 0 {
				block.PreviousBlockHash = HashBlock(blockchain[len(blockchain) - 1])
			} else {
				block.PreviousBlockHash = [32]byte{}
			}
			block.Difficulty = previousBlock.Difficulty * (60 / uint64(previousBlock.MiningTime.Seconds()));
			block.Nonce++
			hashBytes = HashBlock(block)
			hash = binary.BigEndian.Uint64(hashBytes[:])
		}
	}
	block.MiningTime = time.Since(start)
	return block, nil
}
