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
		body := strings.NewReader(fmt.Sprintf("%s%s:%s%s:%s:%s:%s", sender, parametersString, receiver, parametersString, amount, rStr, sStr))
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
		if block.Miner == miner {
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
