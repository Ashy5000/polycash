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
