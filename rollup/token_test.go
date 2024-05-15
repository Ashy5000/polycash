package rollup

import (
	. "cryptocurrency/node_util"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateL2Transaction(t *testing.T) {
	// Arrange
	sender := PublicKey{}
	recipient := PublicKey{}
	amount := uint64(0)
	senderJson, err := json.Marshal(sender)
	if err != nil {
		panic(err)
	}
	recipientJson, err := json.Marshal(recipient)
	if err != nil {
		panic(err)
	}
	expected := "== BEGIN L2 TRANSACTION ==\n" + string(senderJson) + "\n" + string(recipientJson) + "\n" + fmt.Sprint(amount) + "\n"
	// Act
	result, err := CreateL2Transaction(sender, recipient, amount)
	if err != nil {
		panic(err)
	}
	// Assert
	assert.Equal(t, expected, result)
}

func TestCombineL2Transactions(t *testing.T) {
	// Arrange
	transactions := []string{"a\n", "b\n", "c\n"}
	expected := "a\nb\nc\n"
	// Act
	result := CombineL2Transactions(transactions)
	// Assert
	assert.Equal(t, expected, result)
}

func TestSeperateL2Transactions(t *testing.T) {
	// Arrange
	transactionsStr := "== BEGIN L2 TRANSACTION ==\na\nb\n== BEGIN L2 TRANSACTION ==\nc\nd"
	expected := []string{"a\nb\n", "c\nd\n"}
	// Act
	result := SeperateL2Transactions(transactionsStr)
	// Assert
	assert.Equal(t, expected, result)
}

func TestGetL2TokenBalances(t *testing.T) {
	// Arrange
	// Add some transactions to the blockchain
	key := GetKey("../key.json").PublicKey
	keyStr, err := json.Marshal(key)
	if err != nil {
		panic(err)
	}
	body := []byte("== BEGIN L2 TRANSACTION ==\n" + string(keyStr) + "\n" + string(keyStr) + "\n" + "1\n")
	// Sign the body
	privateKey := GetKey("../key.json")
	signature, err := privateKey.X.Sign(body)
	if err != nil {
		panic(err)
	}
	block := Block{
		Transactions: []Transaction{
			{
				Body: body,
				BodySignatures: []Signature{
					{
						S: signature,
					},
				},
			},
		},
	}
	Blockchain = append(Blockchain, block)
	// Act
	result := GetL2TokenBalances()
	// Assert
	assert.Equal(t, map[string]uint64{string(keyStr): 0}, result)
}
