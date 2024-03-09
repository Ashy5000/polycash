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
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	t.Run("It creates a block with valid transaction information", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		senderPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &a,
		}
		recipientPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &b,
		}
		var amount float64
		amount = 123
		blockchain = nil
		transactionHashes[[32]byte{}] = 1
		// Act
		block, err := CreateBlock(senderPublicKey, recipientPublicKey, amount, a, b, [32]byte{})
		if err != nil {
			panic(err)
		}
		// Assert
		assert.Equal(t, senderPublicKey, block.Sender)
		assert.Equal(t, recipientPublicKey, block.Recipient)
		assert.Equal(t, amount, block.Amount)
	})
	t.Run("It creates a block with a valid hash", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		senderPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &a,
		}
		recipientPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &b,
		}
		var amount float64
		amount = 123
		var maxHash uint64
		maxHash = 0x1000000000000000
		transactionHashes[[32]byte{}] = 1
		// Act
		block, err := CreateBlock(senderPublicKey, recipientPublicKey, amount, a, b, [32]byte{})
		if err != nil {
			panic(err)
		}
		// Assert
		hashBytes := HashBlock(block)
		hash := binary.BigEndian.Uint64(hashBytes[:])
		assert.True(t, hash < maxHash)
	})
}
