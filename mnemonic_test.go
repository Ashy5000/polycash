package main

import "testing"

func TestGetMnemonic(t *testing.T) {
	t.Run("It should return a non-empty string", func(t *testing.T) {
		key := GetKey()
		result := GetMnemonic(key)
		if result == "" {
			t.Errorf("GetMnemonic() = %v; want a mnemonic", result)
		}
	})
	t.Run("It should return a reversible mnemonic", func(t *testing.T) {
		key := GetKey()
		mnemonic := GetMnemonic(key)
		result := RestoreMnemonic(mnemonic)
		if result.X.Cmp(key.X) != 0 {
			t.Errorf("RestoreMnemonic(GetMnemonic()) = %v; want %v", result, key)
		}
	})
}
