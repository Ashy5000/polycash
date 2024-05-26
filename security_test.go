// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "cryptocurrency/node_util"
)

func TestApplySecurityLevel(t *testing.T) {
	t.Run("It sets the correct parameters for security level 0", func(t *testing.T) {
		ApplySecurityLevel(0)
		assert.Equal(t, InitialBlockDifficulty, SecurityLevels[0].InitialBlockDifficulty)
		assert.Equal(t, MinimumBlockDifficulty, SecurityLevels[0].MinimumDifficulty)
		assert.Equal(t, BlocksBeforeReward, SecurityLevels[0].BlocksBeforeReward)
	})
	t.Run("It sets the correct parameters for security level 1", func(t *testing.T) {
		ApplySecurityLevel(1)
		assert.Equal(t, InitialBlockDifficulty, SecurityLevels[1].InitialBlockDifficulty)
		assert.Equal(t, MinimumBlockDifficulty, SecurityLevels[1].MinimumDifficulty)
		assert.Equal(t, BlocksBeforeReward, SecurityLevels[1].BlocksBeforeReward)
	})
	t.Run("It sets the correct parameters for security level 2", func(t *testing.T) {
		ApplySecurityLevel(2)
		assert.Equal(t, InitialBlockDifficulty, SecurityLevels[2].InitialBlockDifficulty)
		assert.Equal(t, MinimumBlockDifficulty, SecurityLevels[2].MinimumDifficulty)
		assert.Equal(t, BlocksBeforeReward, SecurityLevels[2].BlocksBeforeReward)
	})
}
