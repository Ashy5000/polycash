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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"
)

// transactionHashes is a map of transaction hashes to their current status. 0 means the transaction is unmined, 1 means the transaction is being mined, and 2 means the transaction has been mined.
var transactionHashes = make(map[[32]byte]int)
var miningTransactions []Transaction

func CreateBlock() (Block, error) {
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
		TimeVerifiers:          []dsa.PublicKey{},
	}
	if block.Difficulty < minimumBlockDifficulty {
		block.Difficulty = minimumBlockDifficulty
	}
	if len(blockchain) > 0 {
		block.PreviousBlockHash = HashBlock(blockchain[len(blockchain)-1])
	} else {
		block.PreviousBlockHash = [32]byte{}
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	fmt.Printf("Mining block with difficulty %d\n", block.Difficulty)
	for hash > maximumUint64/block.Difficulty {
		for i := 0; i < len(miningTransactions); i++ {
			transaction := miningTransactions[i]
			transactionString := fmt.Sprintf("%s:%s:%f:%d", EncodePublicKey(transaction.Sender), EncodePublicKey(transaction.Recipient), transaction.Amount, transaction.Timestamp.UnixNano())
			transactionBytes := []byte(transactionString)
			hash := sha256.Sum256(transactionBytes)
			if transactionHashes[hash] > 1 {
				miningTransactions[i] = miningTransactions[len(miningTransactions)-1]
				miningTransactions = miningTransactions[:len(miningTransactions)-1]
				i--
			}
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
				block.PreviousBlockHash = [32]byte{}
			}
			block.Difficulty = previousBlock.Difficulty * (60 / uint64(previousBlock.MiningTime.Seconds()))
			if block.Difficulty < minimumBlockDifficulty {
				block.Difficulty = minimumBlockDifficulty
			}
			block.Nonce++
			hashBytes = HashBlock(block)
			hash = binary.BigEndian.Uint64(hashBytes[:])
		}
	}
	block.MiningTime = time.Since(start)
	// Convert the block to a string (JSON)
	bodyChars, err := json.Marshal(&block)
	if err != nil {
		panic(err)
	}
	// Ask for time verifiers
	for _, peer := range GetPeers() {
		// Verify that the peer has mined a block (only miners can be time verifiers)
		req, err := http.NewRequest(http.MethodGet, peer+"/identify", nil)
		if err != nil {
			panic(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Peer down.")
			continue
		}
		// Get the response body
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		// Convert the response body to a string
		bodyString := string(bodyBytes)
		// Convert the response body to a big.Int
		peerY, ok := new(big.Int).SetString(bodyString, 10)
		if !ok {
			fmt.Println("Could not convert peer Y to big.Int")
			continue
		}
		// Create a dsa.PublicKey from the big.Int
		peerKey := dsa.PublicKey{
			Y: peerY,
		}
		// Verify that the peer has mined a block
		if IsNewMiner(peerKey, len(blockchain)) {
			fmt.Println("Peer has not mined a block.")
			continue
		}
		// Ask to verify the time
		body := strings.NewReader(string(bodyChars))
		req, err = http.NewRequest(http.MethodGet, peer+"/verifyTime", body)
		if err != nil {
			panic(err)
		}
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Peer down.")
			continue
		}
		// Get the response body
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		if string(bodyBytes) == "invalid" {
			fmt.Println("Time verifier believes block is invalid.")
			continue
		}
		// Unmarshal the response body
		responseBlock := Block{}
		err = json.Unmarshal(bodyBytes, &responseBlock)
		if err != nil {
			panic(err)
		}
		// Set the time verifiers
		block.TimeVerifiers = responseBlock.TimeVerifiers // TODO: Verify time verifiers
	}
	if int64(len(block.TimeVerifiers)) < GetMinerCount(len(blockchain))/5 {
		fmt.Println("Not enough time verifiers.")
		return Block{}, errors.New("lost block")
	}
	return block, nil
}
