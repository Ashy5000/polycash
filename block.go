package main

import (
	"crypto/dsa"
	"math/big"
)

type Block struct {
	Sender    dsa.PublicKey `json:"sender"`
	Recipient dsa.PublicKey `json:"recipient"`
	Miner     dsa.PublicKey `json:"miner"`
	Amount    float64       `json:"amount"`
	Nonce     int64         `json:"nonce"`
	R         big.Int       `json:"R"`
	S         big.Int       `json:"S"`
}
