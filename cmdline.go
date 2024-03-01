// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"crypto/dsa"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
		} else if action == "savestate" {
			// Save the blockchain to a file
			blockchainJson, err := json.Marshal(blockchain)
			if err != nil {
				panic(err)
			}
			err = os.WriteFile("blockchain.json", blockchainJson, 0644)
			if err != nil {
				panic(err)
			}
		} else if action == "loadstate" {
			// Load the blockchain from a file
			blockchainJson, err := os.ReadFile("blockchain.json")
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(blockchainJson, &blockchain)
			if err != nil {
				panic(err)
			}
		} else if action == "exit" {
			return
		} else if action == "help" {
			fmt.Println("Commands:")
			fmt.Println("sync - Sync the blockchain with peers")
			fmt.Println("balance <public key> - Get the balance of a public key")
			fmt.Println("send <public key> <amount> - Send an amount to a public key")
			fmt.Println("keygen - Generate a new key")
			fmt.Println("savestate - Save the blockchain to a file")
			fmt.Println("loadstate - Load the blockchain from a file")
			fmt.Println("exit - Exit the console")
		} else if action == "addpeer" {
			// Get the IP address of the peer
			ip := fields[1]
			// Send a request to the peer server to add the peer
			peerServer := "http://192.168.4.87:8080"
			req, err := http.NewRequest(http.MethodGet, peerServer+"/add_peer/"+ip, nil)
			if err != nil {
				panic(err)
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				panic(err)
			}
			fmt.Println("Peer added successfully!")
		} else if action == "" {
			continue
		} else {
			fmt.Println("Invalid command.")
		}
	}
}
