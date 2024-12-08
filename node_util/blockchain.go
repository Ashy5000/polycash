// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"time"
)

var Blockchain []Block

// GenesisBlock returns a new genesis block.
//
// It returns a Block struct with the following fields:
// - Transactions: a slice of Transaction structs (nil)
// - Miner: a PublicKey struct
// - Nonce: an int64 with value 0
// - MiningTime: an int64 with value 0
// - Difficulty: an uint64 with value 0
// - PreviousBlockHash: a [64]byte array with all elements set to 0
// - Timestamp: a time.Time struct with the zero value
// - TimeVerifierSignatures: a slice of Signature structs (empty)
// - TimeVerifiers: a slice of PublicKey structs (empty)
func GenesisBlock() Block {
	return Block{
		LegacyTransactions:     nil,
		Miner:                  PublicKey{},
		Nonce:                  0,
		MiningTime:             0,
		Difficulty:             0,
		PreviousBlockHash:      [64]byte{},
		Timestamp:              time.Time{},
		TimeVerifierSignatures: []Signature{},
		TimeVerifiers:          []PublicKey{},
	}
}

// Append adds a new block to the blockchain.
//
// It takes a single parameter, `block`, of type `Block`, which represents the block to be appended to the blockchain.
// This function does not return any value.
func Append(block Block) {
	Blockchain = append(Blockchain, block)
}
