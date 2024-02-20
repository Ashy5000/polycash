package main

import (
	"flag"
)

func main() {
	mine := flag.Bool("mine", false, "Set to true to start node as miner")
	serve := flag.Bool("serve", *mine, "Set to true to start node as server")
	flag.Parse()
	if *serve {
		Serve(*mine)
	} else {
		StartCmdLine()
	}
}
