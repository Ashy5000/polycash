// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"crypto/sha256"
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
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%f", senderStr, recipientStr, amount)))
	if transactionHashes[hash] > 0 {
		fmt.Println("No new job. Ignoring mine request.")
		return
	}
	fmt.Println("New job.")
	transactionHashes[hash] = 1
	fmt.Println("Broadcasting job to peers...")
	for _, peer := range GetPeers() {
		// Create a new body
		body := strings.NewReader(string(bodyBytes))
		req, err := http.NewRequest(http.MethodGet, peer+"/mine", body)
		if err != nil {
			panic(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Peer", peer, "is down.")
		}
	}
	if err != nil {
		panic(err)
	}
	rStr := fields[3]
	sStr := fields[4]
	var r big.Int
	var s big.Int
	r.SetString(rStr, 10)
	s.SetString(sStr, 10)
	if !VerifyTransaction(senderKey, recipientKey, strconv.FormatFloat(amount, 'f', -1, 64), r, s) {
		fmt.Println("Transaction is invalid. Ignoring transaction request.")
		return
	}
	block, err := CreateBlock(senderKey, recipientKey, amount, r, s, hash)
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
			fmt.Println("Peer down.")
		}
	}
	fmt.Println("All done!")
}

func HandleBlockRequest(_ http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	block := Block{}
	err = json.Unmarshal(bodyBytes, &block)
	if err != nil {
		panic(err)
	}
	if !VerifyBlock(block) {
		fmt.Println("Block is invalid. Ignoring block request.")
		return
	}
	// Get transaction as string
	transaction := fmt.Sprintf("%s:%s:%f", EncodePublicKey(block.Sender), EncodePublicKey(block.Recipient), block.Amount)
	transactionBytes := []byte(transaction)
	// Get hash of transaction
	hash := sha256.Sum256(transactionBytes)
	fmt.Println("Transaction hash:", hash)
	// Mark transaction as completed
	transactionHashes[hash] = 2
	Append(block)
	fmt.Println("Block appended to local blockchain!")
	if *useLocalPeerList {
		// Broadcast block to peers
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
				fmt.Println("Peer is down.")
			}
		}
	}
	fmt.Println("All done!")
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

func Serve(mine bool, port string) {
	if mine {
		http.HandleFunc("/mine", HandleMineRequest)
	}
	http.HandleFunc("/block", HandleBlockRequest)
	http.HandleFunc("/blockchain", HandleBlockchainRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
