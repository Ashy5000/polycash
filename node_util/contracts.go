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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"os"
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
	// Try with merkle
	contract, ok := GetValue(state.ZenContracts, strconv.FormatUint(c.Location, 10))
	if ok {
		c.Contents = string(contract)
		c.Parties = nil // Zen drops parties from specification (can be replaced by new VM features)
		// By only storing contracts in the merkle tree, the root hash will match between the consensus client and the VM
		// This way, the merkle tree doesn't have to be rebuilt for each transaction
		c.GasUsed = 0
		c.Loaded = true
		return
	}
	// Fallback to legacy
	for _, contract := range state.LegacyContracts {
		if contract.Location == c.Location {
			c.Contents = contract.Contents
			c.Parties = contract.Parties
			c.GasUsed = 0
			c.Loaded = true
			break
		}
	}
}

var executionLocked = false

func (c Contract) Execute(maxGas float64, sender PublicKey) ([]Transaction, StateTransition, float64, error) {
	for executionLocked {
	}
	executionLocked = true
	if !c.Loaded {
		c.LoadContract()
	}
	if !VerifySmartContract(c) {
		Warn("Invalid contract detected.")
		executionLocked = false
		return make([]Transaction, 0), StateTransition{}, 0, nil
	}
	if err := os.WriteFile("contract.blockasm", []byte(c.Contents), 0666); err != nil {
		executionLocked = false
		return nil, StateTransition{}, 0, err
	}
	pendingState := GetPendingState()
	pendingStateSerialized, err := msgpack.Marshal(pendingState)
	err = os.WriteFile("pending_state.msgpack", pendingStateSerialized, 0644)
	if err != nil {
		executionLocked = false
		return nil, StateTransition{}, 0, err
	}
	out, _ := ZkProve([]Contract{c}, []float64{maxGas}, []PublicKey{sender}, CalculateCurrentState())
	scanner := bufio.NewScanner(strings.NewReader(out))
	transactions := make([]Transaction, 0)
	gasUsed := 0.0
	transition := StateTransition{
		LegacyUpdatedData: make(map[string][]byte),
		ZenUpdatedData:    make([]MerkleNode, 0),
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
						executionLocked = false
						return nil, StateTransition{}, 0, err
					}
					fmt.Println("Applying state change:", address, valueBytes)
					if Env.Upgrades.Zen <= len(Blockchain) {
						// Zen insert
						transition.ZenUpdatedData = InsertValue(transition.ZenUpdatedData, address, valueBytes)
					} else {
						// Legacy insert
						transition.LegacyUpdatedData[address] = valueBytes
					}
				}
				continue
			}
			gasUsed, err = strconv.ParseFloat(line[10:], 64)
			if err != nil {
				executionLocked = false
				return nil, StateTransition{}, 0, err
			}
			continue
		}
		words := strings.Split(line, " ")
		var senderY []byte
		err = json.Unmarshal([]byte(words[1]), &senderY)
		if err != nil {
			executionLocked = false
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
			executionLocked = false
			return nil, StateTransition{}, 0, nil
		}
		sender := PublicKey{
			Y: senderY,
		}
		var receiverY []byte
		err = json.Unmarshal([]byte(words[2]), &receiverY)
		if err != nil {
			executionLocked = false
			return nil, StateTransition{}, 0, err
		}
		receiver := PublicKey{
			Y: receiverY,
		}
		subdividedAmount, err := strconv.Atoi(words[3])
		if err != nil {
			executionLocked = false
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
		if Env.Upgrades.Zen <= len(Blockchain) {
			// Zen update
			InsertValue(transition.ZenUpdatedData, strconv.FormatUint(c.Location, 10), []byte(c.Contents))
		} else {
			// Legacy update
			transition.LegacyNewContracts = map[uint64]Contract{
				c.Location: c,
			}
		}
	}
	executionLocked = false
	return transactions, transition, gasUsed, nil
}
