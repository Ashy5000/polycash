package main

import (
	"crypto/dsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
)

func SyncBlockchain() {
	longestLength := 0
	var longestBlockchain []Block
	for _, peer := range GetPeers() {
		res, err := http.Get(fmt.Sprintf("%s/blockchain", peer))
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(res.Body)
		var peerBlockchain []Block
		err = json.Unmarshal(body, &peerBlockchain)
		if err != nil {
			panic(err)
		}
		length := len(peerBlockchain)
		if length > longestLength {
			longestLength = length
			longestBlockchain = peerBlockchain
		}
	}
	blockchain = longestBlockchain
}

func GetBalance(key big.Int) float64 {
	total := 0.0
	for _, block := range blockchain {
		if block.Sender.Y.Cmp(&key) == 0 {
			total -= block.Amount
		} else if block.Recipient.Y.Cmp(&key) == 0 {
			total += block.Amount
		}
	}
	return total
}

func Send(receiver string, amount string) {
	keyJson, err := os.ReadFile("key.json")
	if err != nil {
		panic(err)
	}
	var key dsa.PrivateKey
	err = json.Unmarshal(keyJson, &key)
	sender := key.PublicKey.Y
	if err != nil {
		panic(err)
	}
	parametersString := fmt.Sprintf("&%s&%s&%s", key.PublicKey.Parameters.P, key.PublicKey.Parameters.Q, key.PublicKey.Parameters.G)
	r, s, err := dsa.Sign(rand.Reader, &key, []byte(fmt.Sprintf("%s:%s:%s", sender, receiver, amount)))
	body := strings.NewReader(fmt.Sprintf("%s%s:%s%s:%s:%s:%s", sender, parametersString, receiver, parametersString, amount, r, s))
	for _, peer := range GetPeers() {
		req, err := http.NewRequest(http.MethodGet, peer+"/mine", body)
		if err != nil {
			panic(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
	}
}
