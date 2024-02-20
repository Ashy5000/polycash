package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strings"
)

func main() {
	mine := flag.Bool("mine", false, "Set to true to start node as miner")
	serve := flag.Bool("serve", *mine, "Set to true to start node as server")
	flag.Parse()
	if *serve {
		Serve(*mine)
	} else {
		for {
			fmt.Printf("BlockCMD console: ")
			inputReader := bufio.NewReader(os.Stdin)
			cmd, _ := inputReader.ReadString('\n')
			cmd = cmd[:len(cmd)-1]
			fields := strings.Split(cmd, " ")
			action := fields[0]
			if action == "sync" {
				fmt.Println("Syncing blockchain...")
				longestLength := 0
				var longestBlockchain []Block
				for _, peer := range GetPeers() {
					res, err := http.Get(fmt.Sprintf("%s/blockchain", peer))
					if err != nil {
						panic(err)
					}
					body, err := io.ReadAll(res.Body)
					var peerBlockchain []Block
					err = json.Unmarshal(body, &peerBlockchain)
					if err != nil {
						panic(err)
					}
					length := len(peerBlockchain)
					if length > longestLength {
						longestLength = length
						longestBlockchain = peerBlockchain
					}
				}
				blockchain = longestBlockchain
				fmt.Printf("Blockchain successfully synced!")
				fmt.Printf("Length: %d", longestLength)
			} else if action == "balance" {
				keyStr := fields[1]
				var key big.Int
				key.SetString(keyStr, 10)
				total := 0.0
				for _, block := range blockchain {
					if block.Sender.Y.Cmp(&key) == 0 {
						total -= block.Amount
					} else if block.Recipient.Y.Cmp(&key) == 0 {
						total += block.Amount
					}
				}
				fmt.Println(fmt.Sprintf("Balance of %s: %f", fields[1], total))
			} else if action == "send" {
				sender := fields[1]
				receiver := fields[2]
				amount := fields[3]
				signaturePlaceholder := "&0&0&0"
				body := strings.NewReader(fmt.Sprintf("%s%s:%s%s:%s", sender, signaturePlaceholder, receiver, signaturePlaceholder, amount))
				for _, peer := range GetPeers() {
					req, err := http.NewRequest(http.MethodGet, peer+"/mine", body)
					if err != nil {
						panic(err)
					}
					_, err = http.DefaultClient.Do(req)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
