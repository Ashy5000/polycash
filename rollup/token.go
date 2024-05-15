package rollup

import (
	. "cryptocurrency/node_util"
	"encoding/json"
	"fmt"
	"strings"
)

func CreateL2Transaction(sender PublicKey, recipient PublicKey, amount uint64) string {
	resultStr := ""
	resultStr += "== BEGIN L2 TRANSACTION ==\n"
	senderJson, err := json.Marshal(sender)
	if err != nil {
		panic(err)
	}
	resultStr += string(senderJson) + "\n"
	recipientJson, err := json.Marshal(recipient)
	if err != nil {
		panic(err)
	}
	resultStr += string(recipientJson) + "\n"
	resultStr += fmt.Sprint(amount) + "\n"
	return resultStr
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
