// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
    "crypto/dsa"
    "github.com/stretchr/testify/assert"
    "math/big"
    "testing"
)

func TestDecodePublicKey(t *testing.T) {
    t.Run("It returns a valid dsa.PublicKey when the key is valid", func(t *testing.T) {
        // Act
        key := DecodePublicKey("1&2&3&4")
        // Assert
        assert.NotNil(t, key)
        assert.Equal(t, int64(1), key.Y.Int64())
        assert.Equal(t, int64(2), key.P.Int64())
        assert.Equal(t, int64(3), key.Q.Int64())
        assert.Equal(t, int64(4), key.G.Int64())
    })
}

func TestEncodePublicKey(t *testing.T) {
    t.Run("It returns a valid string when the key is valid", func(t *testing.T) {
        // Arrange
        key := dsa.PublicKey{
            Parameters: dsa.Parameters{
                P: big.NewInt(2),
                Q: big.NewInt(3),
                G: big.NewInt(4),
            },
            Y: big.NewInt(1),
        }
        // Act
        encodedKey := EncodePublicKey(key)
        // Assert
        assert.Equal(t, "1&2&3&4", encodedKey)
    })
}
