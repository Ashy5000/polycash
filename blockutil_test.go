package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMaxMiners(t *testing.T) {
	t.Run("It returns 1 when the length of the blockchain is 0", func(t *testing.T) {
		// Arrange
		blockchain = nil
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(1), maxMiners)
	})
	t.Run("It returns 11 when the length of the blockchain is 50", func(t *testing.T) {
		// Arrange
		blockchain = nil
		for i := 0; i < 50; i++ {
			Append(Block{})
		}
		// Act
		maxMiners := GetMaxMiners()
		// Assert
		assert.Equal(t, int64(11), maxMiners)
	})
}
