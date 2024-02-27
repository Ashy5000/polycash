package main

import (
	"flag"
)

var useLocalPeerList *bool

func main() {
	mine := flag.Bool("mine", false, "Set to true to start node as miner")
	serve := flag.Bool("serve", *mine, "Set to true to start node as server")
	port := flag.String("port", "8080", "Port to listen on (server only)")
	useLocalPeerList = flag.Bool("useLocalPeerList", false, "Set to true to use local peer list")
	flag.Parse()
	if *serve {
		Serve(*mine, *port)
	} else {
		StartCmdLine()
	}
}
