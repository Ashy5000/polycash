// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	. "cryptocurrency/node_interface"
	. "cryptocurrency/node_util"
	. "cryptocurrency/rollup"
	. "cryptocurrency/testing"
	"flag"
	"net/http"
)

func main() {
	mine := flag.Bool("mine", false, "Set to true to start node as miner")
	serve := flag.Bool("serve", *mine, "Set to true to start node as server")
	port := flag.String("port", "8080", "Port to listen on (server only)")
	command := flag.String("command", "exit", "Run a command and exit")
	Verbose = flag.Bool("verbose", false, "Set to true to enable verbose logging")
	benchmark := flag.Bool("benchmark", false, "Set to true to enable benchmarking")
	flag.Parse()
	LoadEnv()
	LoadStateCmd(nil)
	SyncBlockchain(-1)
	if len(Blockchain) == 0 {
		Append(GenesisBlock())
	}
	if *benchmark {
		Benchmark()
		return
	}
	if *mine {
		*serve = true
	}
	if *serve {
		if *mine {
			go Mine()
		}
		http.HandleFunc("/l2Transaction", HandleTransactionRequest)
		Serve(*mine, *port)
	} else {
		if *command == "exit" {
			StartCmdLine()
		} else if *command == "test" {
			StartTest()
		} else {
			RunCmd(*command)
		}
	}
}
