package main

import (
	"crypto/dsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func HandleMineRequest(_ http.ResponseWriter, req *http.Request) {
	fmt.Println("New job, mining...")
	lostBlock = false
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	body := string(bodyBytes)
	fields := strings.Split(body, ":")
	senderStr := fields[0]
	senderKey := DecodePublicKey(senderStr)
	recipientStr := fields[1]
	recipientKey := DecodePublicKey(recipientStr)
	amount, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		panic(err)
	}
	rStr := fields[3]
	sStr := fields[4]
	var r big.Int
	var s big.Int
	r.SetString(rStr, 10)
	s.SetString(sStr, 10)
	isValid := dsa.Verify(&senderKey, []byte(fmt.Sprintf("%s:%s:%s", senderKey.Y, recipientKey.Y, fields[2])), &r, &s)
	if !isValid {
		fmt.Println("Signature is invalid. Ignoring transaction request.")
		return
	}
	block, err := CreateBlock(senderKey, recipientKey, amount)
	if err != nil {
		fmt.Println("Block lost.")
		return
	}
	fmt.Println("Block mined successfully!")
	fmt.Println("Broadcasting block to peers...")
	bodyChars, err := json.Marshal(&block)
	if err != nil {
		panic(err)
	}
	for _, peer := range GetPeers() {
		body := strings.NewReader(string(bodyChars))
		req, err := http.NewRequest(http.MethodGet, peer+"/block", body)
		if err != nil {
			panic(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("All done!")
}

func HandleBlockRequest(_ http.ResponseWriter, req *http.Request) {
	lostBlock = true
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	block := Block{}
	err = json.Unmarshal(bodyBytes, &block)
	if err != nil {
		panic(err)
	}
	hashBytes := HashBlock(block)
	hash := binary.BigEndian.Uint64(hashBytes[:]) // Take the last 64 bits-- we won't ever need more than 64 zeroes.
	if hash > 9223372036854776000/block.Difficulty {
		fmt.Println("Block has invalid hash. Ignoring block request.")
		fmt.Printf("Actual hash: %d\n", hash)
		return
	}
	Append(block)
	fmt.Println("Block appended to local blockchain!")
}

func HandleBlockchainRequest(w http.ResponseWriter, _ *http.Request) {
	blockchainChars, err := json.Marshal(blockchain)
	if err != nil {
		panic(err)
	}
	_, err = io.WriteString(w, string(blockchainChars))
	if err != nil {
		panic(err)
	}
}

func Serve(mine bool) {
	if mine {
		http.HandleFunc("/mine", HandleMineRequest)
	}
	http.HandleFunc("/block", HandleBlockRequest)
	http.HandleFunc("/blockchain", HandleBlockchainRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
