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
