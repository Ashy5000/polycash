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

type Signature struct {
	R big.Int `json:"R"`
	S big.Int `json:"S"`
}

type Block struct {
	Sender                 dsa.PublicKey   `json:"sender"`
	Recipient              dsa.PublicKey   `json:"recipient"`
	Miner                  dsa.PublicKey   `json:"miner"`
	Amount                 float64         `json:"amount"`
	Nonce                  int64           `json:"nonce"`
	SenderSignature        Signature       `json:"signature"`
	MiningTime             time.Duration   `json:"miningTime"`
	Difficulty             uint64          `json:"difficulty"`
	PreviousBlockHash      [32]byte        `json:"previousBlockHash"`
	Timestamp              time.Time       `json:"timestamp"`
	TimeVerifierSignatures []Signature     `json:"timeVerifierSignature"`
	TimeVerifiers          []dsa.PublicKey `json:"timeVerifiers"`
}
