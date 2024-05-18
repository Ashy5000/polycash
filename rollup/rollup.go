package rollup

import (
	"bytes"
	"crypto/sha256"
	. "cryptocurrency/node_util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var nextTransactions = []string{}
var nextTransactionPeerIps = []string{}
var nextTransactionSignatures = [][]byte{}

func HandleTransactionRequest(_ http.ResponseWriter, req *http.Request) {
	fmt.Println("Handling L2 transaction request.")
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	transaction := string(bodyBytes)
	// Add transaction to nextTransactions
	nextTransactions = append(nextTransactions, transaction)
	// Get IP address of requester
	peerIp := req.RemoteAddr
	// Remove port number
	peerIp = strings.Split(peerIp, ":")[0]
	// Add :8080 to IP address
	peerIp += ":8080"
	// Add http:// to IP address
	peerIp = "http://" + peerIp
	fmt.Println(peerIp)
	nextTransactionPeerIps = append(nextTransactionPeerIps, peerIp)
	if len(nextTransactions) >= 1 {
		// Combine transactions
		combinedTransactions := CombineL2Transactions(nextTransactions)
		// Request requesters to sign combined transactions
		for _, peerIp := range nextTransactionPeerIps {
			req, err := http.NewRequest(http.MethodPost, peerIp+"/signL2Transactions", strings.NewReader(combinedTransactions))
			if err != nil {
				panic(err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				Log("Peer is down.", true)
			}
			// Get signature
			signature, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			if bytes.Equal(signature, []byte("invalid")) {
				fmt.Println("Peer sent invalid signature.")
				return // TODO: Remove invalid transactions and request new signatures
			}
			nextTransactionSignatures = append(nextTransactionSignatures, signature)
		}
		fmt.Println("All signatures received.")
		// Create L2 transaction rollup
		rollup := ""
		key := GetKey("")
		keyBytes := EncodePublicKey(key.PublicKey)
		if err != nil {
			panic(err)
		}
		rollup += string(keyBytes)
		rollup += "$"
		rollup += string(keyBytes)
		rollup += "$"
		rollup += "0.0"
		rollup += "$"
		timestamp := time.Now().UnixNano()
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", key.PublicKey.Y, key.PublicKey.Y, "0.0", timestamp)))
		sigBytes, err := key.X.Sign(hash[:])
		if err != nil {
			panic(err)
		}
		signature := Signature{S: sigBytes}
		signatureStr, err := json.Marshal(signature)
		if err != nil {
			panic(err)
		}
		rollup += string(signatureStr)
		rollup += "$"
		rollup += fmt.Sprint(timestamp)
		rollup += "$"
		rollup += "[]"
		rollup += "$"
		bodyStr, err := json.Marshal(combinedTransactions)
		if err != nil {
			panic(err)
		}
		rollup += string(bodyStr)
		rollup += "$"
		signatures := []Signature{}
		for _, signature := range nextTransactionSignatures {
			signature := Signature{S: signature}
			signatures = append(signatures, signature)
		}
		signaturesStr, err := json.Marshal(signatures)
		if err != nil {
			panic(err)
		}
		rollup += string(signaturesStr)
		// Send rollup to all peers
		fmt.Println("Sending rollup to peers...")
		for _, peer := range GetPeers() {
			req, err := http.NewRequest(http.MethodGet, peer+"/mine", strings.NewReader(rollup))
			if err != nil {
				panic(err)
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				Log("Peer is down.", true)
			}
		}
	}
}
