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
	result := CreateL2Transaction(sender, recipient, amount)
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
