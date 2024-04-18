// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"
)

func HashBlock(block Block) [64]byte {
	marshaled, err := json.Marshal(block)
	if err != nil {
		panic(err)
	}
	blockCpy := Block{}
	err = json.Unmarshal(marshaled, &blockCpy)
	if err != nil {
		panic(err)
	}
	blockCpy.MiningTime = time.Minute
	blockCpy.TimeVerifierSignatures = []Signature{}
	blockCpy.TimeVerifiers = []PublicKey{}
	blockCpy.Timestamp = time.Time{}
	for i := range block.Transactions {
		blockCpy.Transactions[i].Timestamp = time.Time{}
	}
	blockBytes := []byte(fmt.Sprintf("%v", blockCpy))
	sum := sha3.Sum512(blockBytes)
	return sum
}
