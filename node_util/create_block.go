// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"time"
)

// TransactionHashes is a map of transaction hashes to their current status. 0 means the transaction is unmined, 1 means the transaction is being mined, and 2 means the transaction has been mined.
var TransactionHashes = make(map[[32]byte]int)
var MiningTransactions []Transaction

func CreateBlock() (Block, error) {
	if len(MiningTransactions) == 0 {
		return Block{}, errors.New("pool dry")
	}
	start := time.Now()
	previousBlock, previousBlockFound := GetLastMinedBlock()
	if !previousBlockFound {
		previousBlock.Difficulty = InitialBlockDifficulty
		previousBlock.MiningTime = time.Minute
	}
	block := Block{
		Miner:                  GetKey("").PublicKey,
		Transactions:           MiningTransactions,
		Nonce:                  0,
		Difficulty:             GetDifficulty(previousBlock.MiningTime, previousBlock.Difficulty),
		Timestamp:              time.Now(),
		TimeVerifierSignatures: []Signature{},
		TimeVerifiers:          []PublicKey{},
		MiningTime:             0,
	}
	if len(Blockchain) > 0 {
		block.PreviousBlockHash = HashBlock(Blockchain[len(Blockchain)-1])
	} else {
		block.PreviousBlockHash = [64]byte{}
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	// Request time verifiers
	block.PreMiningTimeVerifierSignatures, block.PreMiningTimeVerifiers = RequestTimeVerification(block)
	Log(fmt.Sprintf("Mining block with difficulty %d", block.Difficulty), false)
	for hash > MaximumUint64/block.Difficulty {
		i := 0
		for _, transaction := range MiningTransactions {
			transactionString := fmt.Sprintf("%s:%s:%f:%d", EncodePublicKey(transaction.Sender), EncodePublicKey(transaction.Recipient), transaction.Amount, transaction.Timestamp.UnixNano())
			transactionBytes := []byte(transactionString)
			hash := sha256.Sum256(transactionBytes)
			if TransactionHashes[hash] > 1 {
				if i > len(MiningTransactions)-1 {
					Error("Transaction index out of range.", false)
					return Block{}, errors.New("transaction index out of range")
				}
				MiningTransactions[i] = MiningTransactions[len(MiningTransactions)-1]
				MiningTransactions = MiningTransactions[:len(MiningTransactions)-1]
				i--
			}
			i++
		}
		if len(MiningTransactions) > 0 {
			previousBlock, previousBlockFound = GetLastMinedBlock()
			if !previousBlockFound {
				previousBlock.Difficulty = InitialBlockDifficulty
				previousBlock.MiningTime = time.Minute
			}
			if len(Blockchain) > 0 {
				block.PreviousBlockHash = HashBlock(Blockchain[len(Blockchain)-1])
			} else {
				block.PreviousBlockHash = [64]byte{}
			}
			block.Difficulty = GetDifficulty(previousBlock.MiningTime, previousBlock.Difficulty)
			block.Transactions = MiningTransactions
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
	if int64(len(block.TimeVerifiers)) < GetMinerCount(len(Blockchain))/5 {
		Warn("Not enough time verifiers.")
		return Block{}, errors.New("lost block")
	}
	MiningTransactions = []Transaction{}
	return block, nil
}
