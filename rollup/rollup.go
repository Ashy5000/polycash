package rollup

import (
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
var nextTransactionSignatures = []string{}

func HandleTransactionRequest(_ http.ResponseWriter, req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	transaction := string(bodyBytes)
	// Add transaction to nextTransactions
	nextTransactions = append(nextTransactions, transaction)
	// Get IP address of requester
	peerIp := req.RemoteAddr + ":8080"
	nextTransactionPeerIps = append(nextTransactionPeerIps, peerIp)
	if len(nextTransactions) >= 5 {
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
			signatureBytes, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			signature := string(signatureBytes)
			if signature == "invalid" {
				Log("Peer declared transactions to be invalid", true)
				return // TODO: Remove invalid transactions and request new signatures
			}
			nextTransactionSignatures = append(nextTransactionSignatures, signature)
		}
		// Create L2 transaction rollup
		rollup := ""
		key := GetKey("")
		keyBytes, err := json.Marshal(key)
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
		sigStr, err := json.Marshal(sigBytes)
		if err != nil {
			panic(err)
		}
		rollup += string(sigStr)
		rollup += "$"
		rollup += string(timestamp)
		rollup += "$"
		rollup += "[]"
		rollup += "$"
		bodyStr, err := json.Marshal(combinedTransactions)
		if err != nil {
			panic(err)
		}
		rollup += string(bodyStr)
		rollup += "$"
		transactionsStr, err := json.Marshal(nextTransactionSignatures)
		if err != nil {
			panic(err)
		}
		rollup += string(transactionsStr)
		// Send rollup to all peers
		for _, peer := range GetPeers() {
			req, err := http.NewRequest(http.MethodGet, peer+"/mine", strings.NewReader(rollup))
			if err != nil {
				panic(err)
			}
			go SendRequest(req)
		}
	}
}
