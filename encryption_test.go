package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptKey(t *testing.T) {
	t.Run("It encrypts the key.json file", func(t *testing.T) {
		// Act
		EncryptKey("0123456789abcdef")
		// Assert
		assert.True(t, IsKeyEncrypted())
		DecryptKey("0123456789abcdef")
	})
}

func TestDecryptKey(t *testing.T) {
	t.Run("It decrypts the key.json file", func(t *testing.T) {
		// Arrange
		EncryptKey("0123456789abcdef")
		// Act
		DecryptKey("0123456789abcdef")
		// Assert
		assert.False(t, IsKeyEncrypted())
	})
}
