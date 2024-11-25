package node_util

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
)

// Overview
// Zero-knowledge proofs are handled by a Rust program using RiscZero.
// The binary is located at `vm_zk/target/release/host`.
// It handles proof creation and verification.
// It accesses the state of the VM and deployed contracts via a file, merkle.txt.
// Then, it creates a merkle tree from the data.
// It accepts a batch of transactions to create a proof.
// The proof includes the merkle root.

// Go utilities

func WriteZkState(state State) {
	var segments []string
	for location, val := range state.LegacyData {
		var segment string
		segment += location
		segment += ">"
		valHex := hex.EncodeToString(val)
		segment += valHex
		segments = append(segments, segment)
	}
	for _, i := range state.ZenData {
		if i.Data != nil {
			var segment string
			segment += i.Key
			segment += ">"
			valHex := hex.EncodeToString(i.Data)
			segment += valHex
			segments = append(segments, segment)
		}
	}
	result := ""
	for segment := range segments {
		result += segments[segment] + "*"
	}
	if result != "" {
		// Remove last *
		result = result[:len(result)-1]
	}
	err := os.WriteFile("merkle.txt", []byte(result), 0644)
	if err != nil {
		panic(err)
	}
}

func GenerateZkArgs(generate bool, hashes []string, gasLimits []float64, senders []PublicKey) []string {
	if generate {
		// Generate args for a ZK proof
		contractsArg := "contract.blockasm" // Contracts to be included
		hashesArg := ""                     // Contract hashes
		for hash := range hashes {
			hashesArg += hashes[hash] + "%"
		}
		hashesArg = hashesArg[:len(hashesArg)-1]
		gasLimitsArg := "" // Gas limits
		for gasLimit := range gasLimits {
			gasLimitsArg += fmt.Sprintf("%f", gasLimits[gasLimit]) + "%"
		}
		gasLimitsArg = gasLimitsArg[:len(gasLimitsArg)-1]
		sendersArg := "" // Senders
		for sender := range senders {
			sendersArg += hex.EncodeToString(senders[sender].Y) + "%"
		}
		sendersArg = sendersArg[:len(sendersArg)-1]
		merkleArg := "merkle.txt"   // State
		recieptArg := "receipt.bin" // Receipt
		return []string{contractsArg, hashesArg, gasLimitsArg, sendersArg, merkleArg, recieptArg}
	} else {
		// Generate args for a ZK verification
		verificationArg := "V"      // Verification
		receiptArg := "receipt.bin" // Receipt
		return []string{verificationArg, receiptArg}
	}
}

func WriteContractsAggregate(contracts []Contract) {
	var segments []string
	for _, contract := range contracts {
		segments = append(segments, contract.Contents)
	}
	res := ""
	for segment := range segments {
		res += segments[segment] + "*"
	}
	res = res[:len(res)-1]
	err := os.WriteFile("contract.blockasm", []byte(res), 0644)
	if err != nil {
		panic(err)
	}
}

func LoadReceipt() []byte {
	receipt, err := os.ReadFile("receipt.bin")
	if err != nil {
		panic(err)
	}
	return receipt
}

func WriteReceipt(receipt []byte) {
	err := os.WriteFile("receipt.bin", receipt, 0644)
	if err != nil {
		panic(err)
	}
}

func RunZkBinary(args []string) (string, bool) {
	out, err := exec.Command("./vm_zk/target/release/host", args...).Output()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(args)
		return "", false
	}
	return string(out), true
}

func ZkProve(contracts []Contract, gasLimits []float64, senders []PublicKey, state State) (string, []byte) {
	// 1. Write contracts to file
	WriteContractsAggregate(contracts)
	// 2. Write state to file
	WriteZkState(state)
	// 3. Generate arguments
	var hashes []string
	for contract := range contracts {
		contractStr := contracts[contract].Contents
		hash := sha256.Sum256([]byte(contractStr))
		hashes = append(hashes, hex.EncodeToString(hash[:]))
	}
	args := GenerateZkArgs(true, hashes, gasLimits, senders)
	// 4. Run zk binary
	res, ok := RunZkBinary(args)
	if !ok {
		panic("ZK proof failed")
	}
	// 5. Read receipt
	receipt := LoadReceipt()
	return res, receipt
}

func ZKVerify(receipt []byte) bool {
	// 1. Write receipt to file
	WriteReceipt(receipt)
	// 2. Generate arguments
	args := GenerateZkArgs(false, nil, nil, nil)
	// 3. Run zk binary
	_, ok := RunZkBinary(args)
	return ok
}
