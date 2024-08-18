// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"encoding/json"
	"net/http"
	"strings"
)

var NextTransitions = make(map[[32]byte]StateTransition)

// Mine mines a block by creating a new block and broadcasting it to peers.
//
// This function continuously mines blocks by calling the CreateBlock function.
// If there is an error creating the block, the function continues to the next iteration.
// After successfully creating a block, the function logs "Block mined successfully!" and "Broadcasting block to peers...".
// The block is then marshaled into JSON format and broadcasted to all peers.
// The broadcasting process involves sending an HTTP GET request to each peer's "/block" endpoint with the block data.
// If there is an error creating the HTTP request or sending the request, the function logs "Peer down.".
// After broadcasting the block to all peers, the function logs "All done!".
// This process continues indefinitely.
func Mine() {
	for {
		block, err := CreateBlock()
		if err != nil {
			continue
		}
		Log("Block mined successfully!", false)
		Log("Broadcasting block to peers...", true)
		bodyChars, err := json.Marshal(&block)
		if err != nil {
			panic(err)
		}
		for _, peer := range GetPeers() {
			body := strings.NewReader(string(bodyChars))
			req, err := http.NewRequest(http.MethodGet, peer+"/block", body)
			if err != nil {
				panic(err)
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				Log("Peer down.", true)
			}
		}
		Log("All done!", false)
	}
}
