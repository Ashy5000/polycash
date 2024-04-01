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
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	if err != nil {
		panic(err)
	}
	timestamp := fields[5]
	rStr := fields[3]
	sStr := fields[4]
	var r big.Int
	var s big.Int
	r.SetString(rStr, 10)
	s.SetString(sStr, 10)
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%f:%s", senderStr, recipientStr, amount, timestamp)))
	if transactionHashes[hash] > 0 {
		fmt.Println("No new job. Ignoring mine request.")
		return
	}
	fmt.Println("New job.")
	transactionHashes[hash] = 1
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	timestampParsed := time.Unix(0, timestampInt)
	if err != nil {
		panic(err)
	}
	miningTransactions = append(miningTransactions, Transaction{
		Sender:    senderKey,
		Recipient: recipientKey,
		Amount:    amount,
		SenderSignature: Signature{
			R: r,
			S: s,
		},
		Timestamp: timestampParsed,
	})
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
	if !VerifyTransaction(senderKey, recipientKey, strconv.FormatFloat(amount, 'f', -1, 64), r, s) {
		fmt.Println("Transaction is invalid. Ignoring transaction request.")
		return
	}
	block, err := CreateBlock()
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
	for _, transaction := range block.Transactions {
		// Get transaction as string
		transactionString := fmt.Sprintf("%s:%s:%f:%d", EncodePublicKey(transaction.Sender), EncodePublicKey(transaction.Recipient), transaction.Amount, transaction.Timestamp.UnixNano())
		transactionBytes := []byte(transactionString)
		// Get hash of transaction
		hash := sha256.Sum256(transactionBytes)
		fmt.Println("Transaction hash:", hash)
		// Mark transaction as completed
		transactionHashes[hash] = 2
	}
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

func HandleIdentifyRequest(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, GetKey().Y.String())
	if err != nil {
		panic(err)
	}
}

func HandlePeerIpRequest(w http.ResponseWriter, req *http.Request) {
	// Find the IP address of a peer by their public key
	peerKeyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	peerKey := string(peerKeyBytes)
	for _, peer := range GetPeers() {
		req, err := http.NewRequest(http.MethodGet, peer+"/identify", nil)
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Peer is down.")
			continue
		}
		currentPeerKeyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		currentPeerKey := string(currentPeerKeyBytes)
		if currentPeerKey == peerKey {
			_, err := io.WriteString(w, peer)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}

func HandleVerifyTimeRequest(w http.ResponseWriter, req *http.Request) {
	// Verify that the time the block was mined is within a reasonable range of the current time
	// Sign the time with the time verifier's private key
	// This is to prevent miners from mining blocks in the future or the past
	requestBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	request := string(requestBytes)
	// Parse the request (JSON)
	block := Block{}
	err = json.Unmarshal([]byte(request), &block)
	if err != nil {
		panic(err)
	}
	// Get the current time
	currentTime := time.Now()
	// Check if the time the block was mined is within a reasonable range of the current time
	// It cannot be in the future, and it cannot be more than 10 seconds in the past
	if block.Timestamp.After(currentTime) || block.Timestamp.Before(currentTime.Add(-10*time.Second)) {
		_, err := io.WriteString(w, "invalid")
		if err != nil {
			panic(err)
		}
		return
	}
	// Sign the time with the time verifier's (this node's) private key
	key := GetKey()
	r, s, err := dsa.Sign(rand.Reader, &key, []byte(block.Timestamp.String()))
	if err != nil {
		panic(err)
	}
	signature := Signature{
		R: *r,
		S: *s,
	}
	block.TimeVerifierSignatures = append(block.TimeVerifierSignatures, signature)
	block.TimeVerifiers = append(block.TimeVerifiers, key.PublicKey)
	// Convert the block back to JSON
	blockBytes, err := json.Marshal(&block)
	if err != nil {
		panic(err)
	}
	// Send the block back to the requester
	_, err = io.WriteString(w, string(blockBytes))
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
	http.HandleFunc("/identify", HandleIdentifyRequest)
	http.HandleFunc("/peerIp", HandlePeerIpRequest)
	http.HandleFunc("/verifyTime", HandleVerifyTimeRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
