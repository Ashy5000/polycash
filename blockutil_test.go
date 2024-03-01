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
