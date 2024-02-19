package main

import (
	"crypto/dsa"
)

type Block struct {
	Sender    dsa.PublicKey `json:"sender"`
	Recipient dsa.PublicKey `json:"recipient"`
	Miner     dsa.PublicKey `json:"miner"`
	Amount    float64       `json:"amount"`
	Nonce     int64         `json:"nonce"`
}
