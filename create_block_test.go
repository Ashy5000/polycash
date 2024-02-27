package main

import (
	"crypto/dsa"
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	t.Run("It creates a block with valid transaction information", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		senderPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &a,
		}
		recipientPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &b,
		}
		var amount float64
		amount = 123
		// Act
		block, err := CreateBlock(senderPublicKey, recipientPublicKey, amount, a, b)
		if err != nil {
			panic(err)
		}
		// Assert
		assert.Equal(t, senderPublicKey, block.Sender)
		assert.Equal(t, recipientPublicKey, block.Recipient)
		assert.Equal(t, amount, block.Amount)
	})
	t.Run("It creates a block with a valid hash", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		senderPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &a,
		}
		recipientPublicKey := dsa.PublicKey{
			Parameters: dsa.Parameters{},
			Y:          &b,
		}
		var amount float64
		amount = 123
		var maxHash uint64
		maxHash = 0x1000000000000000
		// Act
		block, err := CreateBlock(senderPublicKey, recipientPublicKey, amount, a, b)
		if err != nil {
			panic(err)
		}
		// Assert
		hashBytes := HashBlock(block)
		hash := binary.BigEndian.Uint64(hashBytes[:])
		assert.True(t, hash < maxHash)
	})
}
