package main

import "testing"

func TestGetPeers(t *testing.T) {
	// Test the GetPeers function
	peers := GetPeers()
	if len(peers) == 0 {
		t.Errorf("Expected at least one peer, got none")
	}
}
