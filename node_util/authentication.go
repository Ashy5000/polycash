// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
)

type AuthenticationProof struct {
	PublicKey PublicKey
	Signature Signature
	Data      []byte
}

func SignAuthenticationProof(a *AuthenticationProof) error {
	// Hash the data so the node requesting the signature can't sign arbitrary data
	digest := sha256.Sum256(a.Data)
	// Sign the hash
	key := GetKey("")
	signature, err := key.X.Sign(digest[:])
	if err != nil {
		return err
	}
	a.Signature.S = signature
	return nil
}

func RequestAuthentication(peerIp string) (PublicKey, bool, error) {
	// Generate a random slice of bytes to sign
	// This is to prevent replay attacks
	data := make([]byte, 64)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	// Get the hash of the data
	digest := sha256.Sum256(data)
	// Request the signature from the peer
	req, err := http.NewRequest("GET", peerIp+"/identify", bytes.NewBuffer(data))
	if err != nil {
		return PublicKey{}, false, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return PublicKey{}, false, err
	}
	// Read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return PublicKey{}, false, err
	}
	// Unmarshal the response
	var proof AuthenticationProof
	err = json.Unmarshal(body, &proof)
	if err != nil {
		return PublicKey{}, false, err
	}
	// Verify the signature
	isValid := VerifyAuthenticationProof(&proof, digest[:])
	if !isValid {
		return PublicKey{}, false, nil
	}
	return proof.PublicKey, true, nil
}
