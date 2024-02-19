package main

import (
	"crypto/dsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

func HandleMineRequest(_ http.ResponseWriter, req *http.Request) {
	fmt.Println("New job, mining...")
	lostBlock = false
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	body := string(bodyBytes)
	fields := strings.Split(body, ":")
	var senderY big.Int
	senderY.SetString(fields[0], 10)
	senderKey := dsa.PublicKey{
		Parameters: dsa.Parameters{},
		Y:          &senderY,
	}
	var recipientY big.Int
	recipientY.SetString(fields[1], 10)
	recipientKey := dsa.PublicKey{
		Parameters: dsa.Parameters{},
		Y:          &recipientY,
	}
	amount, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		panic(err)
	}
	block, err := CreateBlock(senderKey, recipientKey, amount)
	if err != nil {
		fmt.Println("Block lost.")
		return
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
			panic(err)
		}
	}
	fmt.Println("All done!")
}

func HandleBlockRequest(_ http.ResponseWriter, req *http.Request) {
	lostBlock = true
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	block := Block{}
	err = json.Unmarshal(bodyBytes, &block)
	if err != nil {
		panic(err)
	}
	Append(block)
	fmt.Println("Block appended to local blockchain!")
}

func Serve() {
	http.HandleFunc("/mine", HandleMineRequest)
	http.HandleFunc("/block", HandleBlockRequest)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
