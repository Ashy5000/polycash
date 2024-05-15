// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerifyTransaction(t *testing.T) {
	t.Run("It should return true if the transaction is valid", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey.Y
		receiver := key.PublicKey.Y
		amount := "0"
		timestamp := time.Now()
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", sender, receiver, amount, timestamp.UnixNano())))
		sig, err := key.X.Sign(hash[:])
		if err != nil {
			panic(err)
		}
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, timestamp, sig)
		assert.True(t, result)
	})
	t.Run("It should return false if the transaction double spends", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey.Y
		receiver := key.PublicKey.Y
		amount := "1"
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender, receiver, amount)))
		sig, err := key.X.Sign(hash[:])
		if err != nil {
			panic(err)
		}
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, time.Now(), sig)
		assert.False(t, result)
	})
	t.Run("It should return false if the transaction signature is invalid", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		amount := "1"
		message := []byte{1, 2, 3, 4}
		sig, err := key.X.Sign(message)
		if err != nil {
			panic(err)
		}
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, time.Now(), sig)
		assert.False(t, result)
	})
}

func TestVerifyMiner(t *testing.T) {
	t.Run("It should return true if the miner is valid", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		result := VerifyMiner(key.PublicKey)
		assert.True(t, result)
	})
	t.Run("It should return false if there are too many miners", func(t *testing.T) {
		key := GetKey("")
		miner := key.PublicKey
		miner.Y = []byte("123")
		Blockchain = nil
		Append(GenesisBlock())
		Append(Block{
			Transactions:      []Transaction{},
			Miner:             miner,
			PreviousBlockHash: HashBlock(GenesisBlock()),
			Difficulty:        1,
		})
		result := VerifyMiner(key.PublicKey)
		assert.False(t, result)
	})
}

func TestVerifyBlock(t *testing.T) {
	t.Run("It should return true if the block is valid", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey
		receiver := key.PublicKey
		amount := 0.0
		timestamp := time.Now()
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", sender.Y, receiver.Y, strconv.FormatFloat(amount, 'f', -1, 64), timestamp.UnixNano())))
		sig, err := key.X.Sign(hash[:])
		if err != nil {
			panic(err)
		}

		miningTransactions = []Transaction{
			{
				Sender:    sender,
				Recipient: receiver,
				Amount:    amount,
				SenderSignature: Signature{
					S: sig,
				},
				Timestamp: timestamp,
			},
		}
		block, err := CreateBlock()
		if err != nil {
			panic(err)
		}
		result := VerifyBlock(block)
		assert.True(t, result)
	})
	t.Run("It should return false if the block has an invalid transaction", func(t *testing.T) {
		key := GetKey("")
		Blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey
		receiver := key.PublicKey
		amount := "1"
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender.Y, receiver.Y, amount)))
		sig, err := key.X.Sign(hash[:])
		if err != nil {
			panic(err)
		}
		block := Block{
			Transactions: []Transaction{
				{
					Sender:          sender,
					Recipient:       receiver,
					Amount:          0,
					SenderSignature: Signature{S: sig},
				},
			},
			PreviousBlockHash: HashBlock(GenesisBlock()),
			Miner:             key.PublicKey,
			Difficulty:        1,
		}
		result := VerifyBlock(block)
		assert.False(t, result)
	})
}
