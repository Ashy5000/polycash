// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ContractParty struct {
	Signature Signature
	PublicKey PublicKey
}

type Contract struct {
	Contents string
	Parties  []ContractParty
	GasUsed  float64
	Location uint64
	Loaded   bool
}

var ExternalStateWriteableValue = []byte("ExternalStateWriteableValue")

func (c Contract) IsNewContract() bool {
	return c.Location == 0
}

func (c Contract) LoadContract() {
	state := CalculateCurrentState()
	for _, contract := range state.Contracts {
		if contract.Location == c.Location {
			c.Contents = contract.Contents
			c.Parties = contract.Parties
			c.GasUsed = contract.GasUsed
			c.Loaded = true
			break
		}
	}
}

func (c Contract) Execute(max_gas float64, sender PublicKey) ([]Transaction, StateTransition, float64, error) {
	if !c.Loaded {
		c.LoadContract()
	}
	if !VerifySmartContract(c) {
		Warn("Invalid contract detected.")
		return make([]Transaction, 0), StateTransition{}, 0, nil
	}
	if err := os.WriteFile("contract.blockasm", []byte(c.Contents), 0666); err != nil {
		return nil, StateTransition{}, 0, err
	}
	contractStr := c.Contents
	hash := sha256.Sum256([]byte(contractStr))
	pendingState := GetPendingState()
	pendingStateSerialized, err := msgpack.Marshal(pendingState)
	err = os.WriteFile("pending_state.msgpack", pendingStateSerialized, 0644)
	if err != nil {
		return nil, StateTransition{}, 0, err
	}
	out, err := exec.Command("./contracts/target/release/contracts", "contract.blockasm", hex.EncodeToString(hash[:]), fmt.Sprintf("%f", max_gas), hex.EncodeToString(sender.Y)).Output()
	if err != nil {
		fmt.Println("Errored with output:", string(out))
		fmt.Println("Error: ", err)
		fmt.Println("Contract hash:", hex.EncodeToString(hash[:]))
		return nil, StateTransition{}, 0, err
	}
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	transactions := make([]Transaction, 0)
	gasUsed := 0.0
	transition := StateTransition{
		UpdatedData: make(map[string][]byte),
	}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}
		if line[:2] != "TX" {
			if len(line) < 10 {
				continue
			}
			if line[:9] != "Gas used:" {
				if line[:14] == "State change: " {
					stateChangeString := line[14:]
					parts := strings.Split(stateChangeString, "|")
					address := parts[0]
					valueHex := parts[1]
					valueBytes, err := hex.DecodeString(valueHex)
					if err != nil {
						return nil, StateTransition{}, 0, err
					}
					transition.UpdatedData[address] = valueBytes
				} else if len(line) >= 25 && line[:24] == "External state change: " {
					stateChangeString := line[24:]
					parts := strings.Split(stateChangeString, "|")
					address := parts[0]
					valueHex := parts[1]
					valueBytes, err := hex.DecodeString(valueHex)
					if err != nil {
						return nil, StateTransition{}, 0, err
					}
					if !bytes.Equal(GetFromState(address), ExternalStateWriteableValue) {
						Warn("Contract attempted to modify external state not marked as writeable.")
						return nil, StateTransition{}, 0, errors.New("contract attempted to modify external state not marked as writeable")
					}
					transition.UpdatedData[address] = valueBytes
				}
				continue
			}
			gasUsed, err = strconv.ParseFloat(line[10:], 64)
			if err != nil {
				return nil, StateTransition{}, 0, err
			}
			continue
		}
		words := strings.Split(line, " ")
		var senderY []byte
		err = json.Unmarshal([]byte(words[1]), &senderY)
		if err != nil {
			return nil, StateTransition{}, 0, err
		}
		senderIsParty := false
		for _, party := range c.Parties {
			if bytes.Equal(party.PublicKey.Y, senderY) {
				senderIsParty = true
				break
			}
		}
		if !senderIsParty {
			Warn("Invalid sender detected.")
			return nil, StateTransition{}, 0, nil
		}
		sender := PublicKey{
			Y: senderY,
		}
		var receiverY []byte
		err = json.Unmarshal([]byte(words[2]), &receiverY)
		if err != nil {
			return nil, StateTransition{}, 0, err
		}
		receiver := PublicKey{
			Y: receiverY,
		}
		subdividedAmount, err := strconv.Atoi(words[3])
		if err != nil {
			return nil, StateTransition{}, 0, err
		}
		amount := float64(subdividedAmount * 1000000)
		transaction := Transaction{
			Sender:            sender,
			Recipient:         receiver,
			Amount:            amount,
			FromSmartContract: true,
		}
		transactions = append(transactions, transaction)
	}
	if c.IsNewContract() {
		transition.NewContracts = map[uint64]Contract{
			c.Location: c,
		}
	}
	return transactions, transition, gasUsed, nil
}
