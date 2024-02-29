package main

import (
	"crypto/dsa"
	"math/big"
	"time"
)

type Block struct {
	Sender            dsa.PublicKey `json:"sender"`
	Recipient         dsa.PublicKey `json:"recipient"`
	Miner             dsa.PublicKey `json:"miner"`
	Amount            float64       `json:"amount"`
	Nonce             int64         `json:"nonce"`
	R                 big.Int       `json:"R"`
	S                 big.Int       `json:"S"`
	MiningTime        time.Duration `json:"miningTime"`
	Difficulty        uint64        `json:"difficulty"`
	PreviousBlockHash [32]byte      `json:"previousBlockHash"`
}
