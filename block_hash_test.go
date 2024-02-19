package main

import (
	"crypto/dsa"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashBlock(t *testing.T) {
	t.Run("It returns a checksum of the block", func(t *testing.T) {
		// Arrange
		block := Block{
			Sender:    dsa.PublicKey{},
			Recipient: dsa.PublicKey{},
			Amount:    0,
			Nonce:     0,
		}
		var emptyHash [32]byte
		// Act
		hash := HashBlock(block)
		// Assert
		assert.NotEqual(t, emptyHash, hash)
	})
}
