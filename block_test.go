// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/json"
	"testing"
	"time"

	. "cryptocurrency/node_util"
	"github.com/stretchr/testify/assert"
)

func TestBlock(t *testing.T) {
	t.Run("It holds the Sender, Recipient, Amount, and Nonce properties", func(t *testing.T) {
		// Arrange
		a := []byte("123")
		b := []byte("321")
		// Act
		block := Block{
			Transactions: []Transaction{
				{
					Sender:    PublicKey{Y: a},
					Recipient: PublicKey{Y: b},
					Amount:    2024,
				},
			},
			Miner:                  PublicKey{},
			Nonce:                  24,
			MiningTime:             0,
			Difficulty:             0,
			PreviousBlockHash:      [64]byte{},
			Timestamp:              time.Time{},
			TimeVerifierSignatures: nil,
			TimeVerifiers:          nil,
		}
		// Assert
		assert.Equal(t, a, block.Transactions[0].Sender.Y)
		assert.Equal(t, b, block.Transactions[0].Recipient.Y)
		assert.Equal(t, float64(2024), block.Transactions[0].Amount)
		assert.Equal(t, int64(24), block.Nonce)
	})
	t.Run("It marshals and unmarshals the block correctly", func(t *testing.T) {
		// Arrange
		a := []byte("123")
		b := []byte("321")
		block := Block{
			Transactions: []Transaction{
				{
					Sender:    PublicKey{Y: a},
					Recipient: PublicKey{Y: b},
					Amount:    2024,
					Timestamp: time.Now(),
				},
			},
			Miner:                  PublicKey{Y: a},
			Nonce:                  24,
			MiningTime:             0,
			Difficulty:             0,
			PreviousBlockHash:      [64]byte{},
			Timestamp:              time.Now(),
			TimeVerifierSignatures: nil,
			TimeVerifiers:          nil,
		}
		marshaled, err := json.Marshal(block)
		if err != nil {
			panic(err)
		}
		unmarshaled := Block{}
		err = json.Unmarshal(marshaled, &unmarshaled)
		if err != nil {
			panic(err)
		}
		timestamp := time.Time{}
		block.Timestamp = timestamp
		for _, transaction := range block.Transactions {
			transaction.Timestamp = timestamp
		}
		// Assert
		assert.Equal(t, HashBlock(block), HashBlock(unmarshaled))
	})
}
