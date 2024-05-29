// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"testing"

	. "cryptocurrency/node_util"
	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	t.Run("It does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Fail(t, "Log() panicked")
			}
		}()
		Log("test", true)
	})
}

func TestWarn(t *testing.T) {
	t.Run("It does not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Fail(t, "Warn() panicked")
			}
		}()
		Warn("test")
	})
}

func TestError(t *testing.T) {
	t.Run("It does not panic when fatal is false", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				assert.Fail(t, "Error() panicked")
			}
		}()
		Error("test", false)
	})
	t.Run("It panics when fatal is true", func(t *testing.T) {
		panicOccurred := false
		defer func() {
			if r := recover(); r != nil {
				panicOccurred = true
			}
		}()
		Error("test", true)
		assert.True(t, panicOccurred)
	})
}
