// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Signature struct {
	R big.Int
	S big.Int
}

func (i Signature) MarshalJSON() ([]byte, error) {
	return []byte(`"` + i.R.String() + "$" + i.S.String() + `"`), nil
}

func (i *Signature) UnmarshalJSON(data []byte) error {
	// Convert data to string
	str := string(data)
	// Remove quotes
	str = strings.Replace(str, `"`, "", -1)
	// Split string into R and S
	parts := strings.Split(str, "$")
	// Convert R and S to big.Int
	i.R.SetString(parts[0], 10)
	i.S.SetString(parts[1], 10)
	return nil
}

type Transaction struct {
	Sender          dsa.PublicKey
	Recipient       dsa.PublicKey
	Amount          float64
	SenderSignature Signature
	Timestamp       time.Time
}

func (i Transaction) MarshalJSON() ([]byte, error) {
	signatureBytes, err := json.Marshal(i.SenderSignature)
	if err != nil {
		panic(err)
	}
	signature := string(signatureBytes)
	result := []byte(EncodePublicKey(i.Sender) + ":" + EncodePublicKey(i.Recipient) + ":" + fmt.Sprintf("%f", i.Amount) + ":" + signature)
	result = []byte(strings.Replace(string(result), `"`, "", -1))
	result = []byte(`"` + string(result) + `"`)
	return result, nil
}

func (i *Transaction) UnmarshalJSON(data []byte) error {
	// Convert data to string
	str := string(data)
	// Remove quotes
	str = strings.Replace(str, `"`, "", -1)
	// Substitute \u0026 with &
	str = strings.Replace(str, "\\u0026", "&", -1)
	// Split string into parts
	parts := strings.Split(str, ":")
	// Convert parts to appropriate types
	i.Sender = DecodePublicKey(parts[0])
	i.Recipient = DecodePublicKey(parts[1])
	i.Amount, _ = strconv.ParseFloat(parts[2], 64)
	var signature Signature
	err := json.Unmarshal([]byte(`"`+parts[3]+`"`), &signature)
	if err != nil {
		fmt.Println(err)
		return err
	}
	i.SenderSignature = signature
	return nil
}

type Block struct {
	Transactions           []Transaction   `json:"transactions"`
	Miner                  dsa.PublicKey   `json:"miner"`
	Nonce                  int64           `json:"nonce"`
	MiningTime             time.Duration   `json:"miningTime"`
	Difficulty             uint64          `json:"difficulty"`
	PreviousBlockHash      [32]byte        `json:"previousBlockHash"`
	Timestamp              time.Time       `json:"timestamp"`
	TimeVerifierSignatures []Signature     `json:"timeVerifierSignature"`
	TimeVerifiers          []dsa.PublicKey `json:"timeVerifiers"`
}
