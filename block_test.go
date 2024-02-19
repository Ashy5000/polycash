package main

import (
	"crypto/dsa"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestBlock(t *testing.T) {
	t.Run("It holds the Sender, Recipient, Amount, and Nonce properties", func(t *testing.T) {
		// Arrange
		var a big.Int
		a.SetUint64(123)
		var b big.Int
		b.SetUint64(321)
		// Act
		block := Block{
			Sender: dsa.PublicKey{
				Parameters: dsa.Parameters{},
				Y:          &a,
			},
			Recipient: dsa.PublicKey{
				Parameters: dsa.Parameters{},
				Y:          &b,
			},
			Amount: 2024,
			Nonce:  24,
		}
		// Assert
		assert.Equal(t, &a, block.Sender.Y)
		assert.Equal(t, &b, block.Recipient.Y)
		assert.Equal(t, float64(2024), block.Amount)
		assert.Equal(t, int64(24), block.Nonce)
	})
}
