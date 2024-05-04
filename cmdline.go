// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/open-quantum-safe/liboqs-go/oqs"
)

var commands = map[string]func([]string){
	"sync":                SyncCmd,
	"balance":             BalanceCmd,
	"send":                SendCmd,
	"deploySmartContract": DeploySmartContractCmd,
	"keygen":              KeygenCmd,
	"showPublicKey":       ShowPublicKeyCmd,
	"encrypt":             EncryptCmd,
	"decrypt":             DecryptCmd,
	"savestate":           SaveStateCmd,
	"loadstate":           LoadStateCmd,
	"addPeer":             AddPeerCmd,
	"bootstrap":           BootstrapCmd,
	"help":                HelpCmd,
	"license":             LicenseCmd,
}

func SyncCmd(fields []string) {
	Log("Syncing blockchain...", false)
	SyncBlockchain()
	Log("Blockchain successfully synced!", false)
	Log(fmt.Sprintf("Length: %d", len(blockchain)), false)
}

func BalanceCmd(fields []string) {
	if len(fields) == 1 {
		publicKey := GetKey().PublicKey.Y
		balance := GetBalance(publicKey)
		fmt.Println(fmt.Sprintf("Balance: %f", balance))
		return
	}
	keyStrFields := fields[1:]
	keyStr := strings.Join(keyStrFields, " ")
	var key []byte
	err := json.Unmarshal([]byte(keyStr), &key)
	if err != nil {
		panic(err)
	}
	balance := GetBalance(key)
	fmt.Println(fmt.Sprintf("Balance: %f", balance))
}

func SendCmd(fields []string) {
	receiverStrFields := fields[1 : len(fields)-1]
	receiverStr := strings.Join(receiverStrFields, " ")
	var receiver []byte
	err := json.Unmarshal([]byte(receiverStr), &receiver)
	if err != nil {
		panic(err)
	}
	amount := fields[len(fields)-1]
	Send(string(receiver), amount)
	Log("Waiting for all workers to finish", true)
	wg.Wait()
	Log("All workers have finished", true)
}

func DeploySmartContractCmd(fields []string) {
	path := fields[1]
	err := DeploySmartContract(path)
	if err != nil {
		panic(err)
	}
	fmt.Println("Smart contract deployed successfully!")
}

func KeygenCmd(fields []string) {
	var privateKey PrivateKey
	sigName := "Dilithium2"
	signer := oqs.Signature{}
	if err := signer.Init(sigName, nil); err != nil {
		Error("Could not initialize Dilithium2 signer", true)
	}
	privateKey.X = signer
	pubKey, err := privateKey.X.GenerateKeyPair()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(privateKey.X.ExportSecretKey()))
	privateKey.PublicKey = PublicKey{
		Y: pubKey,
	}
	keyJson, err := json.Marshal(privateKey)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("key.json", keyJson, 0644)
	if err != nil {
		panic(err)
	}
	// TODO: Implement mnemonics for Dilithium2
	//			mnemonic0 := GetMnemonic(*privateKey.X)
	//			mnemonic1 := GetMnemonic(*privateKey.PublicKey.Y)
	//			mnemonic2 := GetMnemonic(*privateKey.PublicKey.Parameters.P)
	//			mnemonic3 := GetMnemonic(*privateKey.PublicKey.Parameters.Q)
	//			mnemonic4 := GetMnemonic(*privateKey.PublicKey.Parameters.G)
	//			fmt.Println("Mnemonic:")
	//			fmt.Println("Part 0: " + mnemonic0)
	//			fmt.Println("Part 1: " + mnemonic1)
	//			fmt.Println("Part 2: " + mnemonic2)
	//			fmt.Println("Part 3: " + mnemonic3)
	//			fmt.Println("Part 4: " + mnemonic4)
	//			fmt.Println("Write down the mnemonic and keep it safe, or better yet memorize it. It is the ONLY WAY to recover your private key.")
}

func ShowPublicKeyCmd(fields []string) {
	// Show the public key in the key.json file
	publicKey := GetKey().PublicKey
	publicKeyJson, err := json.Marshal(publicKey.Y)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(publicKeyJson))
}

func EncryptCmd(fields []string) {
	// Ask the user for a password
	fmt.Print("Enter a password: ")
	inputReader := bufio.NewReader(os.Stdin)
	password, _ := inputReader.ReadString('\n')
	password = password[:len(password)-1]
	// Encrypt the key
	EncryptKey(password)
}

func DecryptCmd(fields []string) {
	// Ask the user for a password
	fmt.Print("Enter a password: ")
	inputReader := bufio.NewReader(os.Stdin)
	password, _ := inputReader.ReadString('\n')
	password = password[:len(password)-1]
	// Decrypt the key
	DecryptKey(password)
}

func SaveStateCmd(fields []string) {
	blockchainJson, err := json.Marshal(blockchain)
	// Save the blockchain to a file
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("blockchain.json", blockchainJson, 0644)
	if err != nil {
		panic(err)
	}
}

func LoadStateCmd(fields []string) {
	// Load the blockchain from a file
	blockchainJson, err := os.ReadFile("blockchain.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(blockchainJson, &blockchain)
	if err != nil {
		panic(err)
	}
}

func AddPeerCmd(fields []string) {
	if *useLocalPeerList {
		// Add the peer to the local peer list
		AddPeer("http://" + fields[1] + ":8080\n")
		// Add the peer to the peer's peer list
		peerServer := fields[1]
		localIp := fields[2]
		body := strings.NewReader(localIp)
		req, err := http.NewRequest(http.MethodGet, peerServer+"/addPeer", body)
		if err != nil {
			panic(err)
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
	} else {
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
	}
	fmt.Println("Peer added successfully!")
}

func BootstrapCmd(fields []string) {
	Bootstrap()
	fmt.Println("Bootstrap complete!")
}

func HelpCmd(fields []string) {
	fmt.Println("Commands:")
	fmt.Println("sync - Sync the blockchain with peers")
	fmt.Println("balance <public key> - Get the balance of a public key")
	fmt.Println("send <public key> <amount> - Send an amount to a public key")
	fmt.Println("keygen - Generate a new key")
	fmt.Println("savestate - Save the blockchain to a file")
	fmt.Println("loadstate - Load the blockchain from a file")
	fmt.Println("exit - Exit the console")
}

func LicenseCmd(fields []string) {
	license, err := os.ReadFile("COPYING")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(license))
}

func RunCmd(input string) {
	cmds := strings.Split(input, ";")
	for _, cmd := range cmds {
		fields := strings.Split(cmd, " ")
		action := fields[0]
		switch action {
		case "exit":
			os.Exit(0)
		case "":
			return
		default:
			fn, ok := commands[action]
			if ok {
				fn(fields)
			} else {
				fmt.Println("Invalid command")
			}
		}
	}
}

func StartCmdLine() {
	fmt.Println("Copyright (C) 2024 Asher Wrobel")
	fmt.Println("This program comes with ABSOLUTELY NO WARRANTY. This is free software, and you are welcome to redistribute it under certain conditions.")
	fmt.Println("To see the license, type `license`.")
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("BlockCMD console (encrypted: %t): ", IsKeyEncrypted())
		cmd, _ := inputReader.ReadString('\n')
		cmd = cmd[:len(cmd)-1]
		RunCmd(cmd)
	}
}
