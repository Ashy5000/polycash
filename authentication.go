package main

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
	key := GetKey()
	signature, err := key.X.Sign(digest[:])
	if err != nil {
		return err
	}
	a.Signature.S = signature
	return nil
}

func RequestAuthentication(peer_ip string) (PublicKey, bool) {
	// Generate a random slice of bytes to sign
	// This is to prevent replay attacks
	data := make([]byte, 64)
	_, err := rand.Read(data)
	if err != nil {
		panic(err)
	}
	// Request the signature from the peer
	req, err := http.NewRequest("GET", peer_ip+"/sign", bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	// Read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	// Unmarshal the response
	var proof AuthenticationProof
	err = json.Unmarshal(body, &proof)
	// Verify the signature
	isValid := VerifyAuthenticationProof(&proof, data)
	if !isValid {
		return PublicKey{}, false
	}
	return proof.PublicKey, true
}
