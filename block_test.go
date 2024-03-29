// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestBlock(t *testing.T) {
	t.Run("It holds the Sender, Recipient, Amount, and Nonce properties", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		// Act
		block := Block{
			Transactions:           []Transaction{
				{
					Sender: dsa.PublicKey{Y: &a},
					Recipient: dsa.PublicKey{Y: &b},
					Amount: 2024,
				},
			},
			Miner:                  dsa.PublicKey{},
			Nonce:                  24,
			MiningTime:             0,
			Difficulty:             0,
			PreviousBlockHash:      [32]byte{},
			Timestamp:              time.Time{},
			TimeVerifierSignatures: nil,
			TimeVerifiers:          nil,
		}
		// Assert
		assert.Equal(t, &a, block.Transactions[0].Sender.Y)
		assert.Equal(t, &b, block.Transactions[0].Recipient.Y)
		assert.Equal(t, float64(2024), block.Transactions[0].Amount)
		assert.Equal(t, int64(24), block.Nonce)
	})
}
