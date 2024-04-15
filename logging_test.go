package main

import (
	"testing"

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
