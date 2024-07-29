package main

import (
	. "cryptocurrency/node_util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateBlockReward(t *testing.T) {
	t.Run("It returns 1 if there are no miners", func(t *testing.T) {
		// Arrange
		LoadEnv()
		// Act
		reward := CalculateBlockReward(0, 0)
		// Assert
		assert.Equal(t, 1.0, reward)
	})
	t.Run("It returns 0.95 if there is 1 miner", func(t *testing.T) {
		// Arrange
		LoadEnv()
		// Act
		reward := CalculateBlockReward(1, 0)
		// Assert
		assert.Equal(t, 0.95, reward)
	})
	t.Run("It returns 0.99 if there is 1 miner and the Guadalajara update is active", func(t *testing.T) {
		// Arrange
		LoadEnv()
		// Act
		reward := CalculateBlockReward(1, 100)
		// Assert
		assert.Equal(t, 0.99, reward)
	})
}
