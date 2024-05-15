// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"encoding/json"
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
			Log("Peer is down.", true)
			continue
		}
		peerPeersBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		var peerPeers []string
		err = json.Unmarshal(peerPeersBytes, &peerPeers)
		if err != nil {
			panic(err)
		}
		for _, peerPeer := range peerPeers {
			if !PeerKnown(peerPeer) {
				// Add the peer's peers to the list of peers
				AddPeer(peerPeer)
			}
		}
	}
}
