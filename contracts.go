package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
}

func (c Contract) Execute() ([]Transaction, error) {
	fmt.Println("Deploying smart contract.")
	if err := os.WriteFile("contract.blockasm", []byte(c.Contents), 0666); err != nil {
		return nil, err
	}
	out, err := exec.Command("./contracts/target/debug/contracts contract.blockasm").Output()
	if err != nil {
		return nil, err
	}
	fmt.Println("Smart contract output: ", string(out))
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	transactions := make([]Transaction, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 3 {
			continue
		}
		if line[:2] != "TX" {
			continue
		}
		fmt.Println("Transaction initiated by contract: ", line)
		words := strings.Split(line, " ")
		words = append(words[:1], words[2:]...)
		var sender PublicKey
		err = json.Unmarshal([]byte(words[0]), &sender)
		if err != nil {
			return nil, err
		}
		var receiver PublicKey
		err = json.Unmarshal([]byte(words[1]), &receiver)
		if err != nil {
			return nil, err
		}
		subdivided_amount, err := strconv.Atoi(words[2])
		if err != nil {
			return nil, err
		}
		amount := float64(subdivided_amount * 1000000)
		transaction := Transaction{
			Sender:    sender,
			Recipient: receiver,
			Amount:    amount,
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
