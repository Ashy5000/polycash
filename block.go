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
	"strconv"
	"strings"
	"time"
)

type Signature struct {
	S []byte
}

func (i Signature) MarshalJSON() ([]byte, error) {
	sigBytes, err := json.Marshal(i.S)
	if err != nil {
		return []byte(""), err
	}
	sigBytes = []byte(strings.Replace(string(sigBytes), `"`, "-", -1))
	return []byte(`"` + string(sigBytes) + `"`), nil
}

func (i *Signature) UnmarshalJSON(data []byte) error {
	// Convert data to string
	str := string(data)
	// Remove double quotes
	str = strings.Replace(str, `"`, "", -1)
	// Convert dashes to double quotes
	str = strings.Replace(str, "-", `"`, -1)
	// Convert string to byte array
	err := json.Unmarshal([]byte(str), &i.S)
	if err != nil {
		return err
	}
	return nil
}

type Transaction struct {
	Sender          PublicKey
	Recipient       PublicKey
	Amount          float64
	SenderSignature Signature
	Timestamp       time.Time
	Contracts       []Contract
}

func (i Transaction) MarshalJSON() ([]byte, error) {
	signatureBytes, err := json.Marshal(i.SenderSignature)
	if err != nil {
		panic(err)
	}
	signature := string(signatureBytes)
	contractsBytes, err := json.Marshal(i.Contracts)
	if err != nil {
		panic(err)
	}
	contracts := string(contractsBytes)
	result := []byte(EncodePublicKey(i.Sender) + ":" + EncodePublicKey(i.Recipient) + ":" + fmt.Sprintf("%f", i.Amount) + ":" + signature + ":" + strconv.FormatInt(i.Timestamp.UnixNano(), 10) + ":" + contracts)
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
	amount, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return err
	}
	i.Amount = amount
	timestampInt, err := strconv.ParseInt(parts[4], 10, 64)
	if err != nil {
		return err
	}
	i.Timestamp = time.Unix(0, timestampInt)
	var signature Signature
	err = json.Unmarshal([]byte(`"`+parts[3]+`"`), &signature)
	if err != nil {
		return err
	}
	i.SenderSignature = signature
	var contracts []Contract
	err = json.Unmarshal([]byte(`"`+parts[4]+`"`), &contracts)
	if err != nil {
		return err
	}
	i.Contracts = contracts
	return nil
}

type Block struct {
	Transactions                    []Transaction `json:"transactions"`
	Miner                           PublicKey     `json:"miner"`
	Nonce                           int64         `json:"nonce"`
	MiningTime                      time.Duration `json:"miningTime"`
	Difficulty                      uint64        `json:"difficulty"`
	PreviousBlockHash               [64]byte      `json:"previousBlockHash"`
	Timestamp                       time.Time     `json:"timestamp"`
	PreMiningTimeVerifierSignatures []Signature   `json:"preMiningTimeVerifierSignatures"`
	PreMiningTimeVerifiers          []PublicKey   `json:"preMiningTimeVerifiers"`
	TimeVerifierSignatures          []Signature   `json:"timeVerifierSignature"`
	TimeVerifiers                   []PublicKey   `json:"timeVerifiers"`
}
