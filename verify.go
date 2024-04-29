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
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unsafe"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func VerifyTransaction(senderKey PublicKey, recipientKey PublicKey, amount string, timestamp time.Time, sig []byte) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	transactionString := fmt.Sprintf("%s:%s:%s:%d", senderKey.Y, recipientKey.Y, amount, timestamp.UnixNano())
	hash := sha256.Sum256([]byte(transactionString))
	verifier := oqs.Signature{}
	sigName := "Dilithium2"
	if err := verifier.Init(sigName, nil); err != nil {
		Error("Failed to initialize Dilithium2 verifier", true)
	}
	isValid, err := verifier.Verify(hash[:], sig, senderKey.Y)
	if err != nil {
		panic(err)
	}
	if !isValid {
		Warn("Invalid transaction signature detected")
		return false
	}
	if GetBalance(senderKey.Y)-amountFloat < 0 {
		Log("Double-spending detected.", true)
		return false
	}
	return true
}

func VerifyMiner(miner PublicKey) bool {
	if IsNewMiner(miner, len(blockchain)) && GetMinerCount(len(blockchain)) >= GetMaxMiners() {
		Log(fmt.Sprintf("Miner count: %d", GetMinerCount(len(blockchain))), true)
		Log(fmt.Sprintf("Maximum miner count: %d", GetMaxMiners()), true)
		return false
	}
	return true
}

func VerifyTransactions(transactions []Transaction) bool {
	for _, transaction := range transactions {
		if !VerifyTransaction(transaction.Sender, transaction.Recipient, strconv.FormatFloat(transaction.Amount, 'f', -1, 64), transaction.Timestamp, transaction.SenderSignature.S) {
			Log("Block has invalid transaction/transaction signature. Ignoring block request.", true)
			return false
		}
	}
	return true
}

func DetectFork(block Block) bool {
	for _, b := range blockchain {
		if b.PreviousBlockHash == block.PreviousBlockHash {
			Warn("Block creates a fork.")
			Log("The node software is designed to handle this edge case, so operations can continue as normal.", false)
			Log("This is most likely a result of latency between miners. If the issue persists, the network may be under attack or a bug may be present; please open an issue on the GitHub repository.", true)
			Log("The blockchain will be re-synced to stay on the longest chain.", true)
			SyncBlockchain()
			return true
		}
	}
	return false
}

func DetectDuplicateBlock(hashBytes [64]byte) bool {
	isDuplicate := false
	for _, b := range blockchain {
		if HashBlock(b) == hashBytes {
			isDuplicate = true
		}
	}
	return isDuplicate
}

func VerifyBlock(block Block) bool {
	isValid := true
	isValid = unsafe.Sizeof(block) <= uintptr(maxBlockSize) && isValid
	isValid = VerifyTransactions(block.Transactions) && isValid
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	isValid = hash <= maximumUint64/block.Difficulty && isValid
	isValid = !DetectDuplicateBlock(hashBytes) && isValid
	isValid = !DetectFork(block) && isValid
	if len(blockchain) > 0 && block.PreviousBlockHash != HashBlock(blockchain[len(blockchain)-1]) {
		Log("Block has invalid previous block hash. Ignoring block request.", true)
		Log("The block could be on a different fork.", true)
		Log("The blockchain will be re-synced to stay on the longest chain.", true)
		SyncBlockchain()
		isValid = false
	}
	isValid = VerifyMiner(block.Miner) && isValid
	// Get the correct difficulty for the block
	lastMinedBlock, found := GetLastMinedBlock()
	if !found {
		lastMinedBlock.Difficulty = initialBlockDifficulty
		lastMinedBlock.MiningTime = time.Minute
	}
	correctDifficulty := GetDifficulty(lastMinedBlock.MiningTime, lastMinedBlock.Difficulty)
	if block.Difficulty != correctDifficulty || block.Difficulty < minimumBlockDifficulty {
		Warn("Invalid difficulty detected.")
		Log("The node software is designed to prevent difficulty manipulation, so this invalid difficulty will not cause issues for the network.", false)
		Log(fmt.Sprintf("Expected difficulty: %d", correctDifficulty), true)
		Log(fmt.Sprintf("Actual difficulty: %d", block.Difficulty), true)
		isValid = false
	}
	if block.Timestamp.After(time.Now()) {
		Log("Block has invalid timestamp. Ignoring block request.", true)
		Log("Timestamp is in the future.", true)
		isValid = false
	}
	if !VerifyTimeVerifiers(block, block.TimeVerifiers, block.TimeVerifierSignatures, false) || !VerifyTimeVerifiers(block, block.PreMiningTimeVerifiers, block.PreMiningTimeVerifierSignatures, true) {
		Log("Block has invalid time verifiers. Ignoring block request.", true)
		isValid = false
	}
	return isValid
}

func VerifyTimeVerifiers(block Block, verifiers []PublicKey, signatures []Signature, premining bool) bool {
	if len(verifiers) != len(signatures) {
		Log("Signature count does not match verifier count.", true)
		return false
	}
	oqsVerifier := oqs.Signature{}
	sigName := "Dilithium2"
	if err := oqsVerifier.Init(sigName, nil); err != nil {
		Error("Failed to initialize Dilithium2 verifier", true)
	}
	if premining {
		for i, verifier := range verifiers {
			valid, err := oqsVerifier.Verify([]byte(fmt.Sprintf("%d", block.Timestamp.UnixNano())), signatures[i].S, verifier.Y)
			if err != nil {
				panic(err)
			}
			if !valid {
				Warn("Invalid time verifier signature detected")
				return false
			}
		}
	} else {
		for i, verifier := range verifiers {
			valid, err := oqsVerifier.Verify([]byte(fmt.Sprintf("%d", block.Timestamp.Add(block.MiningTime).UnixNano())), signatures[i].S, verifier.Y)
			if err != nil {
				panic(err)
			}
			if !valid {
				Warn("Invalid time verifier signature detected")
				return false
			}
		}
	}
	// Ensure all verifiers are unique
	verifierMap := make(map[string]bool)
	for _, verifier := range verifiers {
		if verifierMap[string(verifier.Y)] {
			Log("Time verifier is not unique.", true)
			return false
		}
		verifierMap[string(verifier.Y)] = true
	}
	// Ensure all verifiers are miners
	for _, verifier := range verifiers {
		if IsNewMiner(verifier, len(blockchain)+1) {
			Log("Time verifier is not a miner.", true)
			return false
		}
	}
	// Ensure there are enough verifiers
	if len(verifiers) < GetMinVerifiers() {
		Log("Not enough time verifiers.", true)
		return false
	}
	return true
}

func GetMinVerifiers() int {
	// Get the last block
	lastBlock := blockchain[len(blockchain)-1]
	// Get the number of verifiers in the last block
	lastVerifierCount := len(lastBlock.TimeVerifiers)
	// Get the minimum number of verifiers
	minVerifiers := int(float64(lastVerifierCount) * 0.66)
	return minVerifiers
}

func VerifySmartContract(contract Contract) bool {
	contractStr, err := json.Marshal(contract)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(contractStr)
	for _, party := range contract.Parties {
		verifier := oqs.Signature{}
		sigName := "Dilithium2"
		if err := verifier.Init(sigName, nil); err != nil {
			Error("Failed to initialize Dilithium2 verifier", true)
		}
		isValid, err := verifier.Verify(hash[:], party.Signature.S, party.PublicKey.Y)
		if err != nil {
			panic(err)
		}
		if !isValid {
			Warn("Invalid smart contract signature detected.")
			return false
		}
	}
	return true
}
