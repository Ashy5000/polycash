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
	"fmt"
	"math/big"
	"strconv"
)

func VerifyTransaction(senderKey dsa.PublicKey, recipientKey dsa.PublicKey, amount string, r big.Int, s big.Int) bool {
	amountFloat, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	isValid := dsa.Verify(&senderKey, []byte(fmt.Sprintf("%s:%s:%s", senderKey.Y, recipientKey.Y, strconv.FormatFloat(amountFloat, 'f', -1, 64))), &r, &s)
	return isValid
}

func VerifyMiner(miner dsa.PublicKey) bool {
	if IsNewMiner(miner, len(blockchain)) && GetMinerCount() >= GetMaxMiners() {
		println("Miner count: ", GetMinerCount())
		println("Maximum miner count: ", GetMaxMiners())
		return false
	}
	return true
}

func VerifyBlock(block Block) bool {
	if !VerifyTransaction(block.Sender, block.Recipient, strconv.FormatFloat(block.Amount, 'f', -1, 64), block.R, block.S) {
		return false
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	if hash > 9223372036854776000/block.Difficulty {
		fmt.Println("Block has invalid hash. Ignoring block request.")
		fmt.Printf("Actual hash: %d\n", hash)
		return false
	}
	if !VerifyMiner(block.Miner) {
		return false
	}
	return true
}
