package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplySecurityLevel(t *testing.T) {
	t.Run("It sets the correct parameters for security level 0", func(t *testing.T) {
		ApplySecurityLevel(0)
		assert.Equal(t, initialBlockDifficulty, securityLevels[0].InitialBlockDifficulty)
		assert.Equal(t, minimumBlockDifficulty, securityLevels[0].MinimumDifficulty)
		assert.Equal(t, blocksBeforeSpendable, securityLevels[0].BlocksBeforeSpendable)
	})
	t.Run("It sets the correct parameters for security level 1", func(t *testing.T) {
		ApplySecurityLevel(1)
		assert.Equal(t, initialBlockDifficulty, securityLevels[1].InitialBlockDifficulty)
		assert.Equal(t, minimumBlockDifficulty, securityLevels[1].MinimumDifficulty)
		assert.Equal(t, blocksBeforeSpendable, securityLevels[1].BlocksBeforeSpendable)
	})
	t.Run("It sets the correct parameters for security level 2", func(t *testing.T) {
		ApplySecurityLevel(2)
		assert.Equal(t, initialBlockDifficulty, securityLevels[2].InitialBlockDifficulty)
		assert.Equal(t, minimumBlockDifficulty, securityLevels[2].MinimumDifficulty)
		assert.Equal(t, blocksBeforeSpendable, securityLevels[2].BlocksBeforeSpendable)
	})
}
