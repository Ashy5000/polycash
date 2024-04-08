// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Mine() {
	for {
		block, err := CreateBlock()
		if err != nil {
			continue
		}
		fmt.Println("Block mined successfully!")
		fmt.Println("Broadcasting block to peers...")
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
				fmt.Println("Peer down.")
			}
		}
		fmt.Println("All done!")
	}
}
