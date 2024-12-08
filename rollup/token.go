package rollup

import (
	. "cryptocurrency/node_util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

func CreateL2Transaction(sender PublicKey, recipient PublicKey, amount uint64) (string, error) {
	resultStr := ""
	resultStr += "== BEGIN L2 TRANSACTION ==\n"
	senderJson, err := json.Marshal(sender)
	if err != nil {
		return "", err
	}
	resultStr += string(senderJson) + "\n"
	recipientJson, err := json.Marshal(recipient)
	if err != nil {
		return "", err
	}
	resultStr += string(recipientJson) + "\n"
	resultStr += fmt.Sprint(amount) + "\n"
	return resultStr, nil
}

func CombineL2Transactions(transactions []string) string {
	resultStr := ""
	for _, transaction := range transactions {
		resultStr += transaction
	}
	return resultStr
}

func SeperateL2Transactions(transactionsStr string) []string {
	lines := strings.Split(transactionsStr, "\n")
	var transactions []string
	for _, line := range lines {
		if line == "== BEGIN L2 TRANSACTION ==" {
			transactions = append(transactions, "")
			continue
		}
		transactions[len(transactions)-1] += line + "\n"
	}
	return transactions
}

func BodyContainsL2Transactions(body string) bool {
	return strings.HasPrefix(body, "== BEGIN L2 TRANSACTION ==")
}

func GetL2TokenBalances() map[string]uint64 {
	balances := make(map[string]uint64)
	for _, block := range Blockchain {
		for _, transaction := range ExtractTransactions(block) {
			body := transaction.Body
			if !BodyContainsL2Transactions(string(body)) {
				continue
			}
			l2Transactions := SeperateL2Transactions(string(body))
			for _, l2Transaction := range l2Transactions {
				lines := strings.Split(l2Transaction, "\n")
				sender := PublicKey{}
				recipient := PublicKey{}
				amount := uint64(0)
				err := json.Unmarshal([]byte(lines[0]), &sender)
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal([]byte(lines[1]), &recipient)
				if err != nil {
					panic(err)
				}
				amount, err = strconv.ParseUint(lines[2], 10, 64)
				if err != nil {
					panic(err)
				}
				// Verify the body signature
				verifier := oqs.Signature{}
				sigName := "Dilithium3"
				err = verifier.Init(sigName, nil)
				if err != nil {
					panic(err)
				}
				foundValidSignature := false
				for _, signature := range transaction.BodySignatures {
					// Check if the signature is valid using the sender's public key
					isValid, err := verifier.Verify(transaction.Body, signature.S, sender.Y)
					if err != nil {
						panic(err)
					}
					if isValid {
						foundValidSignature = true
						break
					}
				}
				if !foundValidSignature {
					continue
				}
				if _, ok := balances[lines[0]]; !ok {
					balances[lines[0]] = 0
				}
				if _, ok := balances[lines[1]]; !ok {
					balances[lines[1]] = 0
				}
				if balances[lines[0]] < amount {
					continue
				}
				balances[lines[0]] -= amount
				balances[lines[1]] += amount
			}
		}
	}
	return balances
}
