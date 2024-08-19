// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"net/http"
	"testing"

	. "cryptocurrency/node_util"
	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	t.Run("It returns a PrivateKey when the key.json file is found", func(t *testing.T) {
		// Act
		key := GetKey("")
		// Assert
		assert.NotNil(t, key)
	})
}

func TestSyncBlockchain(t *testing.T) {
	t.Run("It sets the blockchain to the longest blockchain from the peers or panics", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		LoadEnv()
		// Act
		defer func() {
			// Assert
			if r := recover(); r == nil {
				assert.NotNil(t, Blockchain)
			}
		}()
		SyncBlockchain(-1)
	})
}

func TestGetBalance(t *testing.T) {
	t.Run("It returns 0 when the blockchain is empty", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		var key []byte
		// Act
		balance := GetBalance(key)
		// Assert
		assert.Equal(t, float64(0), balance)
	})
	t.Run("It returns the correct balance of a key", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		Append(GenesisBlock())
		key := []byte("123")
		sender := PublicKey{
			Y: []byte("321"),
		}
		receiver := PublicKey{
			Y: key,
		}
		Append(Block{
			Transactions: []Transaction{
				{
					Sender:    sender,
					Recipient: receiver,
					Amount:    100,
				},
			},
			Miner: sender,
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
		Wg.Add(1)
		assert.NotPanics(t, func() {
			SendRequest(req)
		})
	})
}

func TestGetLastMinedBlock(t *testing.T) {
	t.Run("It returns the last block mined by the key", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		Append(GenesisBlock())
		key := GetKey("").PublicKey
		block := Block{
			Miner: key,
		}
		Append(block)
		// Act
		lastMinedBlock, found := GetLastMinedBlock(GetKey("").PublicKey.Y)
		// Assert
		assert.True(t, found)
		assert.Equal(t, block, lastMinedBlock)
	})
	t.Run("It returns false when the key has not mined any blocks", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		Append(GenesisBlock())
		// Act
		_, found := GetLastMinedBlock(GetKey("").PublicKey.Y)
		// Assert
		assert.False(t, found)
	})
}

func TestGetMaxMiners(t *testing.T) {
	t.Run("It returns 1 when the length of the blockchain is 0", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(1), maxMiners)
	})
	t.Run("It returns 2 when the length of the blockchain is 40", func(t *testing.T) {
		// Arrange
		Blockchain = nil
		for i := 0; i < 40; i++ {
			Append(Block{})
		}
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(2), maxMiners)
	})
}
