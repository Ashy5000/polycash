// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/binary"
	"testing"

	. "cryptocurrency/node_util"

	"github.com/stretchr/testify/assert"
)

func TestCreateBlock(t *testing.T) {
	t.Run("It creates a block with valid transaction information", func(t *testing.T) {
		// Arrange
		a := []byte("123")
		b := []byte("321")
		senderPublicKey := PublicKey{
			Y: a,
		}
		recipientPublicKey := PublicKey{
			Y: b,
		}
		var amount float64
		amount = 123
		Blockchain = nil
		TransactionHashes[[32]byte{}] = 1
		// Act
		MiningTransactions = []Transaction{
			{
				Sender:    senderPublicKey,
				Recipient: recipientPublicKey,
				Amount:    amount,
			},
		}
		block, err := CreateBlock()
		if err != nil {
			panic(err)
		}
		// Assert
		assert.Equal(t, senderPublicKey, block.Transactions[0].Sender)
		assert.Equal(t, recipientPublicKey, block.Transactions[0].Recipient)
		assert.Equal(t, amount, block.Transactions[0].Amount)
	})
	t.Run("It creates a block with a valid hash", func(t *testing.T) {
		// Arrange
		a := []byte("123")
		b := []byte("321")
		senderPublicKey := PublicKey{
			Y: a,
		}
		recipientPublicKey := PublicKey{
			Y: b,
		}
		var amount float64
		amount = 123
		var maxHash uint64
		maxHash = 0x1000000000000000
		TransactionHashes[[32]byte{}] = 1
		// Act
		MiningTransactions = []Transaction{
			{
				Sender:    senderPublicKey,
				Recipient: recipientPublicKey,
				Amount:    amount,
			},
		}
		block, err := CreateBlock()
		if err != nil {
			panic(err)
		}
		// Assert
		hashBytes := HashBlock(block, 0)
		hash := binary.BigEndian.Uint64(hashBytes[:])
		assert.True(t, hash < maxHash)
	})
}
