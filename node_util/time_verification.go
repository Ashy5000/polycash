// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func RequestTimeVerification(block Block) ([]Signature, []PublicKey) {
	Log("Requesting time verification", true)
	var signatures []Signature
	var publicKeys []PublicKey
	// Convert the block to a string (JSON)
	bodyChars, err := json.Marshal(&block)
	if err != nil {
		panic(err)
	}
	for _, peer := range GetPeers() {
		// Get the peer's public key
		peerKey, validSig := RequestAuthentication(peer)
		if !validSig {
			Log("Peer has invalid signature.", true)
			continue
		}
		// Verify that the peer has mined a block
		if IsNewMiner(peerKey, len(blockchain)+1) {
			Log("Peer has not mined a block.", true)
			continue
		}
		// Ask to verify the time
		body := strings.NewReader(string(bodyChars))
		req, err := http.NewRequest(http.MethodGet, peer+"/verifyTime", body)
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
		var publicKey PublicKey
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
