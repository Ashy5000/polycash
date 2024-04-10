package main

import (
	"crypto/dsa"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strings"
)

func RequestTimeVerification(block Block) ([]Signature, []dsa.PublicKey) {
  Log("Requesting time verification", true)
  var signatures []Signature
  var publicKeys []dsa.PublicKey
	// Convert the block to a string (JSON)
	bodyChars, err := json.Marshal(&block)
	if err != nil {
		panic(err)
	}
	for _, peer := range GetPeers() {
    if int64(len(block.TimeVerifiers)) >= GetMinerCount(len(blockchain))/5 {
			break
		}
		// Verify that the peer has mined a block (only miners can be time verifiers)
		req, err := http.NewRequest(http.MethodGet, peer+"/identify", nil)
		if err != nil {
			panic(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			Log("Peer down.", true)
			continue
		}
		// Get the response body
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		// Convert the response body to a string
		bodyString := string(bodyBytes)
		// Convert the response body to a big.Int
		peerY, ok := new(big.Int).SetString(bodyString, 10)
		if !ok {
			Log("Could not convert peer Y to big.Int", true)
			continue
		}
		// Create a dsa.PublicKey from the big.Int
		peerKey := dsa.PublicKey{
			Y: peerY,
		}
		// Verify that the peer has mined a block
		if IsNewMiner(peerKey, len(blockchain) + 1) {
			Log("Peer has not mined a block.", true)
			continue
		}
		// Ask to verify the time
		body := strings.NewReader(string(bodyChars))
		req, err = http.NewRequest(http.MethodGet, peer+"/verifyTime", body)
		if err != nil {
			panic(err)
		}
		res, err = http.DefaultClient.Do(req)
		if err != nil {
			Log("Peer down.", true)
			continue
		}
		// Get the response body
		bodyBytes, err = io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		if string(bodyBytes) == "invalid" {
		  Warn("verifier believes block is invalid.")
			continue
		}
		// Split the response body into the signature and the public key
		split := strings.Split(string(bodyBytes), "%")
		// Unmarshal the signature
		var signature Signature
		err = json.Unmarshal([]byte(split[0]), &signature)
		if err != nil {
			panic(err)
		}
		// Unmarshal the public key
		var publicKey dsa.PublicKey
		err = json.Unmarshal([]byte(split[1]), &publicKey)
		if err != nil {
			panic(err)
		}
		// Add the time verifier to the block
		publicKeys = append(publicKeys, publicKey)
		// Add the time verifier signature to the block
		signatures = append(signatures, signature)
    Log("Got verification.", true)
	}
  return signatures, publicKeys
}
