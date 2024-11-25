package rollup

import (
	"bytes"
	. "cryptocurrency/node_util"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

var pendingTransactions []string
var listening = false

func SendL2Transaction(sender PublicKey, recipient PublicKey, amount uint64) {
	transaction, err := CreateL2Transaction(sender, recipient, amount)
	if err != nil {
		panic(err)
	}
	// Send transaction to peers
	for _, peer := range GetPeers() {
		// Send transaction to peer
		req, err := http.NewRequest("POST", peer+"/l2Transaction", bytes.NewBuffer([]byte(transaction)))
		if err != nil {
			panic(err)
		}
		Wg.Add(1)
		go SendRequest(req)
	}
	// Listen for signing requests if not already listening
	if !listening {
		http.HandleFunc("/signL2Transaction", HandleSignL2TransactionRequest)
		go func() {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				panic(err)
			}
		}()
		listening = true
	}
	// Add transaction to pending transactions
	pendingTransactions = append(pendingTransactions, transaction)
}

func HandleSignL2TransactionRequest(w http.ResponseWriter, r *http.Request) {
	// Get transaction
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	// Split into transactions
	transactions := SeperateL2Transactions(string(body))
	myTransactionsCount := 0
	for _, transaction := range transactions {
		// Get sender (2nd line)
		senderStr := strings.Split(transaction, "\n")[1]
		sender := PublicKey{}
		err := json.Unmarshal([]byte(senderStr), &sender)
		if err != nil {
			panic(err)
		}
		if bytes.Equal(sender.Y, GetKey("").PublicKey.Y) {
			myTransactionsCount++
		}
		// Ensure transaction is in pending transactions
		found := false
		for i, pendingTransaction := range pendingTransactions {
			if pendingTransaction == transaction {
				found = true
				// Remove transaction from pending transactions
				pendingTransactions = append(pendingTransactions[:i], pendingTransactions[i+1:]...)
				break
			}
		}
		if !found {
			_, err := w.Write([]byte("invalid"))
			if err != nil {
				panic(err)
			}
		}
	}
	if myTransactionsCount != len(pendingTransactions) {
		_, err := w.Write([]byte("invalid"))
		if err != nil {
			panic(err)
		}
	}
	// Sign combined transactions
	key := GetKey("")
	signature, err := key.X.Sign(body)
	if err != nil {
		panic(err)
	}
	// Send signature
	marshaledSignature, err := json.Marshal(signature)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(marshaledSignature)
	if err != nil {
		panic(err)
	}
}
