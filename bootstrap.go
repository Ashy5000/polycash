package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Bootstrap() {
	// Connect to all peers' peers
	peers := GetPeers()
	for _, peer := range peers {
		// Get the peer's peers
		req, err := http.NewRequest(http.MethodGet, peer+"/peers", nil)
		if err != nil {
			panic(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Peer is down.")
			continue
		}
		peerPeersBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var peerPeers []string
		err = json.Unmarshal(peerPeersBytes, &peerPeers)
		for _, peerPeer := range peerPeers {
			if !PeerKnown(peerPeer) {
				// Add the peer's peers to the list of peers
				AddPeer(peerPeer)
			}
		}
	}
}
