// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"math/big"
	"time"
)

var blockchain []Block

func GenesisBlock() Block {
	return Block{
		Sender:            dsa.PublicKey{},
		Recipient:         dsa.PublicKey{},
		Miner:             dsa.PublicKey{},
		Amount:            0,
		Nonce:             0,
		R:                 big.Int{},
		S:                 big.Int{},
		MiningTime:        0,
		Difficulty:        0,
		PreviousBlockHash: [32]byte{},
		Timestamp:         time.Time{},
	}
}

func Append(block Block) {
	blockchain = append(blockchain, block)
}
