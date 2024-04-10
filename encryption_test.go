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
