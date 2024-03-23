// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

func VerifyTransaction(senderKey dsa.PublicKey, recipientKey dsa.PublicKey, amount string, r big.Int, s big.Int) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", senderKey.Y, recipientKey.Y, amount)))
	isValid := dsa.Verify(&senderKey, hash[:], &r, &s)
	if !isValid {
		return false
	}
	if GetBalance(*senderKey.Y)-amountFloat < 0 {
		return false
	}
	return true
}

func VerifyMiner(miner dsa.PublicKey) bool {
	if IsNewMiner(miner, len(blockchain)) && GetMinerCount() >= GetMaxMiners() {
		println("Miner count: ", GetMinerCount())
		println("Maximum miner count: ", GetMaxMiners())
		return false
	}
	return true
}

func VerifyBlock(block Block) bool {
	if !VerifyTransaction(block.Sender, block.Recipient, strconv.FormatFloat(block.Amount, 'f', -1, 64), block.R, block.S) {
		fmt.Println("Block has invalid transaction/transaction signature. Ignoring block request.")
		return false
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	if hash > maximumUint64/block.Difficulty {
		fmt.Println("Block has invalid hash. Ignoring block request.")
		fmt.Printf("Actual hash: %d\n", hash)
		return false
	}
	for _, b := range blockchain {
		if HashBlock(b) == hashBytes {
			fmt.Println("Block already exists. Ignoring block request.")
			return false
		}
		if b.PreviousBlockHash == block.PreviousBlockHash {
			fmt.Println("Block creates a fork.")
			fmt.Println("This is most likely a result of latency between miners. If the issue persists, the network may be under attack or a bug may be present; please open an issue on the GitHub repository.")
			fmt.Println("The blockchain will be re-synced to stay on the longest chain.")
			SyncBlockchain()
			return true
		}
	}
	if len(blockchain) > 0 && block.PreviousBlockHash != HashBlock(blockchain[len(blockchain)-1]) {
		fmt.Println("Block has invalid previous block hash. Ignoring block request.")
		fmt.Println("The block could be on a different fork.")
		fmt.Println("The blockchain will be re-synced to stay on the longest chain.")
		SyncBlockchain()
		return false
	}
	if !VerifyMiner(block.Miner) {
		return false
	}
	// Get the correct difficulty for the block
	lastMinedBlock := Block{
		Difficulty: initialBlockDifficulty,
		MiningTime: time.Minute,
	}
	if len(blockchain) > 0 {
		isGenesis := true
		for _, b := range blockchain {
			if isGenesis {
				isGenesis = false
				continue
			}
			if b.Miner.Y.Cmp(block.Miner.Y) == 0 {
				lastMinedBlock = b
			}
		}
	}
	correctDifficulty := GetDifficulty(lastMinedBlock.MiningTime, lastMinedBlock.Difficulty)
	if block.Difficulty != correctDifficulty {
		fmt.Println("Block has invalid difficulty. Ignoring block request.")
		fmt.Println("Expected difficulty: ", correctDifficulty)
		fmt.Println("Actual difficulty: ", block.Difficulty)
		return false
	}
	return true
}
