package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func AddPeer(ip string) {
	f, err := os.OpenFile("peers.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	if _, err = f.WriteString(ip); err != nil {
		panic(err)
	}
}

func GetPeers() []string {
	if *useLocalPeerList {
		file, err := os.Open("peers.txt")
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		scanner := bufio.NewScanner(file)

		var result []string
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		return result
	} else {
		peerServer := "http://192.168.4.87:8080"
		// Send a request to the peer server to get the list of peers
		res, err := http.Get(peerServer + "/get_peers")
		if err != nil {
			panic(err)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		// Split response on newline
		result := strings.Split(string(body), "\n")
		// Remove empty strings
		for i := 0; i < len(result); i++ {
			if result[i] == "" {
				result = append(result[:i], result[i+1:]...)
				i--
			}
		}
		// Add http:// to each peer
		// Add :8080 to each peer (port 8080 is used for the peer server, and other ports can be used for local peer lists)
		for i, peer := range result {
			result[i] = "http://" + peer + ":8080"
			fmt.Println("Peer: " + result[i])
		}
		return result
	}
}
