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
	"net/http"
	"testing"
)

func TestGetKey(t *testing.T) {
	t.Run("It returns a dsa.PrivateKey when the key.json file is found", func(t *testing.T) {
		// Act
		key := GetKey()
		// Assert
		assert.NotNil(t, key)
	})
}

func TestSyncBlockchain(t *testing.T) {
	t.Run("It sets the blockchain to the longest blockchain from the peers or panics", func(t *testing.T) {
		// Arrange
		blockchain = nil
		// Act
		defer func() {
			// Assert
			if r := recover(); r == nil {
				assert.NotNil(t, blockchain)
			}
		}()
		SyncBlockchain()
	})
}

func TestGetBalance(t *testing.T) {
	t.Run("It returns 0 when the blockchain is empty", func(t *testing.T) {
		// Arrange
		blockchain = nil
		var key big.Int
		key.SetString("1234567890", 10)
		// Act
		balance := GetBalance(key)
		// Assert
		assert.Equal(t, float64(0), balance)
	})
	t.Run("It returns the correct balance of a key", func(t *testing.T) {
		// Arrange
		blockchain = nil
		Append(GenesisBlock())
		var key big.Int
		key.SetString("1234567890", 10)
		sender := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          big.NewInt(987654321),
		}
		receiver := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &key,
		}
		Append(Block{
			Sender:    sender,
			Recipient: receiver,
			Miner:     sender,
			Amount:    100,
		})
		// Act
		balance := GetBalance(key)
		// Assert
		assert.Equal(t, float64(100), balance)
	})
}

func TestSendRequest(t *testing.T) {
	t.Run("It does not panic when the request is successful", func(t *testing.T) {
		// Arrange
		req := &http.Request{}
		// Act & Assert
		assert.NotPanics(t, func() {
			SendRequest(req)
		})
	})
}

func TestGetMaxMiners(t *testing.T) {
	t.Run("It returns 1 when the length of the blockchain is 0", func(t *testing.T) {
		// Arrange
		blockchain = nil
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(1), maxMiners)
	})
	t.Run("It returns 11 when the length of the blockchain is 50", func(t *testing.T) {
		// Arrange
		blockchain = nil
		for i := 0; i < 50; i++ {
			Append(Block{})
		}
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(11), maxMiners)
	})
}
