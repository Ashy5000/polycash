package main

import (
	"crypto/dsa"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDecodePublicKey(t *testing.T) {
	t.Run("It returns a valid dsa.PublicKey when the key is valid", func(t *testing.T) {
		// Act
		key := DecodePublicKey("1&2&3&4")
		// Assert
		assert.NotNil(t, key)
		assert.Equal(t, int64(1), key.Y.Int64())
		assert.Equal(t, int64(2), key.P.Int64())
		assert.Equal(t, int64(3), key.Q.Int64())
		assert.Equal(t, int64(4), key.G.Int64())
	})
}

func TestEncodePublicKey(t *testing.T) {
	t.Run("It returns a valid string when the key is valid", func(t *testing.T) {
		// Arrange
		key := dsa.PublicKey{
			Parameters: dsa.Parameters{
				P: big.NewInt(2),
				Q: big.NewInt(3),
				G: big.NewInt(4),
			},
			Y: big.NewInt(1),
		}
		// Act
		encodedKey := EncodePublicKey(key)
		// Assert
		assert.Equal(t, "1&2&3&4", encodedKey)
	})
}
