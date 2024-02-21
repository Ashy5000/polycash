package main

import (
	"bufio"
	"crypto/dsa"
	"crypto/rand"
	"encoding/json"
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
			receiver := fields[1]
			amount := fields[2]
			Send(receiver, amount)
		} else if action == "keygen" {
			var privateKey dsa.PrivateKey
			var params dsa.Parameters
			err := dsa.GenerateParameters(&params, rand.Reader, dsa.ParameterSizes(0))
			if err != nil {
				panic(err)
			}
			privateKey.Parameters = params
			err = dsa.GenerateKey(&privateKey, rand.Reader)
			if err != nil {
				panic(err)
			}
			keyJson, err := json.Marshal(privateKey)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("key.json", keyJson, 0644)
		}
	}
}
