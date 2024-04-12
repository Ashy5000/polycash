package main

import (
	"encoding/json"
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"strings"
)

type PublicKey struct {
	Y []byte
}

type PrivateKey struct {
	PublicKey PublicKey
	X oqs.Signature
}

func (i PrivateKey) MarshalJSON() ([]byte, error) {
	pubKey, err := json.Marshal(i.PublicKey)
	if err != nil {
		return []byte(""), err
	}
	privKey, err := json.Marshal(i.X.ExportSecretKey())
	if err != nil {
		return []byte(""), err
	}
	result := []byte(string(pubKey) + "-" + string(privKey))
	result = []byte(strings.Replace(string(result), `"`, "'", -1))
	result = []byte(`"` + string(result) + `"`)
	return result, nil
}

func (i *PrivateKey) UnmarshalJSON(data []byte) error {
	// Convert data to string
	str := string(data)
	// Split string into parts
	parts := strings.Split(str, "-")
	
	var pubKey PublicKey
	pubKeyStr := parts[0]
	pubKeyStr = strings.Replace(pubKeyStr, "'", `"`, -1)
	pubKeyStr = pubKeyStr[1:]
	err := json.Unmarshal([]byte(pubKeyStr), &pubKey)
	if err != nil {
		return err
	}
	var privKeyBytes []byte
	privKeyStr := parts[1]
	privKeyStr = privKeyStr[:len(privKeyStr)-1]
	privKeyStr = strings.Replace(privKeyStr, "'", `"`, -1)
	err = json.Unmarshal([]byte(privKeyStr), &privKeyBytes)
	privKey := oqs.Signature{}
	sigName := "Dilithium2"
	if err := privKey.Init(sigName, privKeyBytes); err != nil {
		panic(err)
	}
	
	i.PublicKey = pubKey
	i.X = privKey
	
	return nil
}