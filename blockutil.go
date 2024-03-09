// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/dsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetKey() dsa.PrivateKey {
	keyJson, err := os.ReadFile("key.json")
	if err != nil {
		panic(err)
	}
	var key dsa.PrivateKey
	err = json.Unmarshal(keyJson, &key)
	if err != nil {
		panic(err)
	}
	return key
}

func SyncBlockchain() {
	longestLength := 0
	var longestBlockchain []Block
	errCount := 0
	for _, peer := range GetPeers() {
		res, err := http.Get(fmt.Sprintf("%s/blockchain", peer))
		if err != nil {
			errCount++
			continue
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
	if errCount >= len(GetPeers()) {
		panic("Could not sync blockchain. All peers down.")
	}
	fmt.Println("Blockchain successfully synced!")
	fmt.Printf("%d out of %d peers responded.\n", len(GetPeers())-errCount, len(GetPeers()))
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
		if block.Miner.Y.Cmp(&key) == 0 {
			total++
		}
	}
	return total
}

func SendRequest(req *http.Request) {
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
}

func Send(receiver string, amount string) {
	key := GetKey()
	sender := key.PublicKey.Y
	parametersString := fmt.Sprintf("&%s&%s&%s", key.PublicKey.Parameters.P, key.PublicKey.Parameters.Q, key.PublicKey.Parameters.G)
	r, s, err := dsa.Sign(rand.Reader, &key, []byte(fmt.Sprintf("%s:%s:%s", sender, receiver, amount)))
	if err != nil {
		panic(err)
	}
	rStr := r.String()
	sStr := s.String()
	for _, peer := range GetPeers() {
		body := strings.NewReader(fmt.Sprintf("%s%s:%s%s:%s:%s:%s:%d", sender, parametersString, receiver, parametersString, amount, rStr, sStr, time.Now().UnixNano()))
		req, err := http.NewRequest(http.MethodGet, peer+"/mine", body)
		if err != nil {
			panic(err)
		}
		go SendRequest(req)
	}
}

func GetLastMinedBlock() (Block, bool) {
	pubKey := GetKey().PublicKey.Y
	for i := len(blockchain) - 1; i >= 0; i-- {
		block := blockchain[i]
		if block.Miner.Y.Cmp(pubKey) == 0 {
			return block, true
		}
	}
	return Block{}, false
}

func IsNewMiner(miner dsa.PublicKey, maxBlockPosition int) bool {
	i := 0
	for _, block := range blockchain {
		if i > maxBlockPosition {
			break
		}
		if block.Miner.Y.Cmp(miner.Y) == 0 {
			return false
		}
		i++
	}
	return true
}

func GetMinerCount() int64 {
	var result int64
	result = 0
	i := 0
	for _, block := range blockchain {
		if IsNewMiner(block.Miner, i-1) {
			result++
		}
		i++
	}
	return result
}

func GetMaxMiners() int64 {
	x := float64(len(blockchain))
	if x == 0 {
		return 1
	}
	t := 8.2
	maxRes := max(x/2, 0)
	top := math.Round(maxRes/t) * 10 * x
	resFloat := (top / (3 * x)) + 1
	res := int64(math.Round(resFloat))
	return res
}
