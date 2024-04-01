package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func Mine() {
	for {
		block, err := CreateBlock()
		if err != nil {
			continue
		}
		fmt.Println("Block mined successfully!")
		fmt.Println("Broadcasting block to peers...")
		bodyChars, err := json.Marshal(&block)
		if err != nil {
			panic(err)
		}
		for _, peer := range GetPeers() {
			body := strings.NewReader(string(bodyChars))
			req, err := http.NewRequest(http.MethodGet, peer+"/block", body)
			if err != nil {
				panic(err)
			}
			_, err = http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Peer down.")
			}
		}
		fmt.Println("All done!")
	}
}
