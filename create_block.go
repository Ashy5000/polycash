// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// transactionHashes is a map of transaction hashes to their current status. 0 means the transaction is unmined, 1 means the transaction is being mined, and 2 means the transaction has been mined.
var transactionHashes = make(map[[32]byte]int)
var miningTransactions []Transaction

func CreateBlock() (Block, error) {
	if len(miningTransactions) == 0 {
		return Block{}, errors.New("pool dry")
	}
	start := time.Now()
	previousBlock, previousBlockFound := GetLastMinedBlock()
	if !previousBlockFound {
		previousBlock.Difficulty = initialBlockDifficulty
		previousBlock.MiningTime = time.Minute
	}
	block := Block{
		Miner:                  GetKey().PublicKey,
		Transactions:           miningTransactions,
		Nonce:                  0,
		Difficulty:             GetDifficulty(previousBlock.MiningTime, previousBlock.Difficulty),
		Timestamp:              time.Now(),
		TimeVerifierSignatures: []Signature{},
		TimeVerifiers:          []PublicKey{},
		MiningTime:             0,
	}
	if len(blockchain) > 0 {
		block.PreviousBlockHash = HashBlock(blockchain[len(blockchain)-1])
	} else {
		block.PreviousBlockHash = [64]byte{}
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	// Request time verifiers
	block.PreMiningTimeVerifierSignatures, block.PreMiningTimeVerifiers = RequestTimeVerification(block)
	Log(fmt.Sprintf("Mining block with difficulty %d", block.Difficulty), false)
	for hash > maximumUint64/block.Difficulty {
		i := 0
		for _, transaction := range miningTransactions {
			transactionString := fmt.Sprintf("%s:%s:%f:%d", EncodePublicKey(transaction.Sender), EncodePublicKey(transaction.Recipient), transaction.Amount, transaction.Timestamp.UnixNano())
			transactionBytes := []byte(transactionString)
			hash := sha256.Sum256(transactionBytes)
			if transactionHashes[hash] > 1 {
				if i > len(miningTransactions)-1 {
					Error("Transaction index out of range.", false)
					return Block{}, errors.New("transaction index out of range")
				}
				miningTransactions[i] = miningTransactions[len(miningTransactions)-1]
				miningTransactions = miningTransactions[:len(miningTransactions)-1]
				i--
			}
			i++
		}
		if len(miningTransactions) > 0 {
			previousBlock, previousBlockFound = GetLastMinedBlock()
			if !previousBlockFound {
				previousBlock.Difficulty = initialBlockDifficulty
				previousBlock.MiningTime = time.Minute
			}
			if len(blockchain) > 0 {
				block.PreviousBlockHash = HashBlock(blockchain[len(blockchain)-1])
			} else {
				block.PreviousBlockHash = [64]byte{}
			}
			block.Difficulty = GetDifficulty(previousBlock.MiningTime, previousBlock.Difficulty)
			block.Transactions = miningTransactions
			block.Nonce++
			hashBytes = HashBlock(block)
			hash = binary.BigEndian.Uint64(hashBytes[:])
		} else {
			Log("Pool dry.", false)
			return Block{}, errors.New("pool dry")
		}
	}
	block.MiningTime = time.Since(start)
	// Ask for time verifiers
	block.TimeVerifierSignatures, block.TimeVerifiers = RequestTimeVerification(block)
	if int64(len(block.TimeVerifiers)) < GetMinerCount(len(blockchain))/5 {
		Warn("Not enough time verifiers.")
		return Block{}, errors.New("lost block")
	}
	miningTransactions = []Transaction{}
	return block, nil
}
