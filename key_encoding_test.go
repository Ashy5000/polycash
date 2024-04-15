// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodePublicKey(t *testing.T) {
	t.Run("It returns a non-nil PublicKey when the key is valid", func(t *testing.T) {
		// Act
		key := DecodePublicKey("1234")
		// Assert
		assert.NotNil(t, key)
	})
	t.Run("It returns a PublicKey that, when encoded, is the same as the original key", func(t *testing.T) {
		// Arrange
		originalKey := "1234"
		// Act
		key := DecodePublicKey(originalKey)
		// Assert
		assert.Equal(t, "[210]", EncodePublicKey(key))
	})
}

func TestEncodePublicKey(t *testing.T) {
	t.Run("It returns a valid string when the key is valid", func(t *testing.T) {
		// Arrange
		key := PublicKey{
			Y: []byte("1234"),
		}
		// Act
		encodedKey := EncodePublicKey(key)
		// Assert
		assert.Equal(t, "[49 50 51 52]", encodedKey)
	})
	t.Run("It returns a string that, when decoded, is the same as the original key", func(t *testing.T) {
		// Arrange
		originalKey := "1234"
		// Act
		key := DecodePublicKey(originalKey)
		// Assert
		assert.Equal(t, "[210]", EncodePublicKey(key))
	})
}
