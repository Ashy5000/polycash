package main

import (
	"crypto/dsa"
	"math/big"
	"strings"
)

func DecodePublicKey(keyString string) dsa.PublicKey {
	fields := strings.Split(keyString, "&")
	var y big.Int
	y.SetString(fields[0], 10)
	var p big.Int
	p.SetString(fields[1], 10)
	var q big.Int
	q.SetString(fields[2], 10)
	var g big.Int
	g.SetString(fields[3], 10)
	publicKey := dsa.PublicKey{
		Parameters: dsa.Parameters{
			P: &p,
			Q: &q,
			G: &g,
		},
		Y: &y,
	}
	return publicKey
}
