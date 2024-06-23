package testing

import (
	"fmt"
	"os"
	"os/exec"
)

func MoveBlockchain(start bool) {
	var from string
	var to string
	if start {
		from = "blockchain.json"
		to = "blockchain_moved.json"
	} else {
		from = "blockchain_moved.json"
		to = "blockchain.json"
	}
	err := os.Rename(from, to)
	if err != nil {
		panic(err)
	}
}

func StartNode(port int) {
	fmt.Println("Starting node...")
	cmd := exec.Command("./builds/node/node_linux-amd64", "-serve", "-mine", "-port", fmt.Sprintf("%d", port))
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
}

func StartNodes() {
	ports := []int{8081}
	for _, port := range ports {
		go StartNode(port)
	}
}

func StartNetwork() {
	MoveBlockchain(true)
	StartNodes()
}

func StopNetwork() {
	MoveBlockchain(false)
}
