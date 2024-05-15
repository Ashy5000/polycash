package rollup

import (
	. "cryptocurrency/node_util"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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
	transactions := []string{}
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
		for _, transaction := range block.Transactions {
			body := transaction.Body
			if !BodyContainsL2Transactions(string(body)) {
				continue
			}
			l2_transactions := SeperateL2Transactions(string(body))
			for _, l2_transaction := range l2_transactions {
				lines := strings.Split(l2_transaction, "\n")
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
