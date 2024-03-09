package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenesisBlock(t *testing.T) {
	t.Run("GenesisBlock() returns a block with a previousBlockHash of all zeroes", func(t *testing.T) {
		// Act
		block := GenesisBlock()
		// Assert
		assert.Equal(t, [32]byte{}, block.PreviousBlockHash)
	})
}

func TestAppend(t *testing.T) {
	t.Run("Append() appends a block to the blockchain", func(t *testing.T) {
		// Arrange
		block := Block{}
		// Act
		Append(block)
		// Assert
		assert.Equal(t, 1, len(blockchain))
	})
}
