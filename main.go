// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"flag"
)

var useLocalPeerList = &[]bool{true}[0]

func main() {
	mine := flag.Bool("mine", false, "Set to true to start node as miner")
	serve := flag.Bool("serve", *mine, "Set to true to start node as server")
	port := flag.String("port", "8080", "Port to listen on (server only)")
	useLocalPeerList = flag.Bool("useLocalPeerList", true, "Set to true to use local peer list and fully decentralize (slower, but more secure)")
	flag.Parse()
	Append(GenesisBlock())
	if *serve {
		Serve(*mine, *port)
	} else {
		StartCmdLine()
	}
}
