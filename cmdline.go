package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func StartCmdLine() {
	for {
		fmt.Printf("BlockCMD console: ")
		inputReader := bufio.NewReader(os.Stdin)
		cmd, _ := inputReader.ReadString('\n')
		cmd = cmd[:len(cmd)-1]
		fields := strings.Split(cmd, " ")
		action := fields[0]
		if action == "sync" {
			fmt.Println("Syncing blockchain...")
			SyncBlockchain()
			fmt.Printf("Blockchain successfully synced!")
			fmt.Printf("Length: %d", len(blockchain))
		} else if action == "balance" {
			keyStr := fields[1]
			var key big.Int
			key.SetString(keyStr, 10)
			balance := GetBalance(key)
			fmt.Println(fmt.Sprintf("Balance of %s: %f", fields[1], balance))
		} else if action == "send" {
			sender := fields[1]
			receiver := fields[2]
			amount := fields[3]
			Send(sender, receiver, amount)
		}
	}
}
