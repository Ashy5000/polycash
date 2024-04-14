// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func GetKey() PrivateKey {
	keyJson, err := os.ReadFile("key.json")
	if err != nil {
		panic(err)
	}
	var key PrivateKey
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
	Log("Blockchain successfully synced!", false)
	Log(fmt.Sprintf("%d out of %d peers responded.", len(GetPeers())-errCount, len(GetPeers())), false)
	blockchain = longestBlockchain
}

func GetBalance(key []byte) float64 {
	total := 0.0
	miningTotal := 0.0
	isGenesis := true
	for i, block := range blockchain {
		if isGenesis {
			isGenesis = false
			continue
		}
		for _, transaction := range block.Transactions {
			if bytes.Equal(transaction.Sender.Y, key) {
				total -= transaction.Amount
				if len(blockchain) > 50 {
					fee := 0.01
					total -= fee
				}
			} else if bytes.Equal(transaction.Recipient.Y, key) {
				total += transaction.Amount
			}
		}
		if bytes.Equal(block.Miner.Y, key) {
			miningTotal += 1.0
			lastBlock := blockchain[i-1]
			miningTotal += float64(len(block.TimeVerifiers)-len(lastBlock.TimeVerifiers)) * 0.1
			if len(blockchain) > 50 {
				fees := 0.0
				for _, _ = range block.Transactions {
					fees += 0.01
				}
			}
		}
	}
	if int(miningTotal) > blocksBeforeSpendable { // A miner must mine n blocks before they can spend their mining rewards
		total += miningTotal
	}
	return total
}

func SendRequest(req *http.Request) {
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		wg.Done()
		return
	}
	wg.Done()
}

func Send(receiver string, amount string) {
	key := GetKey()
	sender := key.PublicKey.Y
	timestamp := time.Now().UnixNano()
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", sender, receiver, amount, timestamp)))
	sigBytes, err := key.X.Sign(hash[:])
	sig := Signature{
		S: sigBytes,
	}
	if err != nil {
		panic(err)
	}
	sigStr, err := json.Marshal(sig)
	senderStr := EncodePublicKey(key.PublicKey)
	receiverStr := EncodePublicKey(PublicKey{Y: []byte(receiver)})
	for _, peer := range GetPeers() {
		Log("Sending transaction to peer: "+peer, false)
		body := strings.NewReader(fmt.Sprintf("%s:%s:%s:%s:%d", senderStr, receiverStr, amount, sigStr, timestamp))
		req, err := http.NewRequest(http.MethodGet, peer+"/mine", body)
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go SendRequest(req)
	}
}

func GetLastMinedBlock() (Block, bool) {
	pubKey := GetKey().PublicKey.Y
	for i := len(blockchain) - 1; i > 0; i-- {
		block := blockchain[i]
		if bytes.Equal(block.Miner.Y, pubKey) {
			return block, true
		}
	}
	return Block{}, false
}

func IsNewMiner(miner PublicKey, maxBlockPosition int) bool {
	isGenesis := true
	for i, block := range blockchain {
		if isGenesis {
			isGenesis = false
			continue
		}
		if i > maxBlockPosition {
			break
		}
		if bytes.Equal(block.Miner.Y, miner.Y) {
			return false
		}
	}
	return true
}

func GetMinerCount(maxBlockPosition int) int64 {
	var result int64
	result = 0
	isGenesis := true
	for i, block := range blockchain {
		if i > maxBlockPosition {
			break
		}
		if isGenesis {
			isGenesis = false
			continue
		}
		if IsNewMiner(block.Miner, i-1) {
			result++
		}
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
