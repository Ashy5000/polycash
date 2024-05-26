// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package node_util

import (
	"bufio"
	"encoding/json"
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
	if *UseLocalPeerList {
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
			Log("Peer: "+result[i], true)
		}
		return result
	}
}

func PeerKnown(ip string) bool {
	peers := GetPeers()
	for _, peer := range peers {
		if peer == ip {
			return true
		}
	}
	return false
}

func ConnectToPeer(ip string) {
	// Add peer to list
	AddPeer(ip)
	// Get IP address of self
	type IP struct {
		Query string
	}
	ipReq, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		panic(err)
	}
	defer ipReq.Body.Close()

	body, err := io.ReadAll(ipReq.Body)
	if err != nil {
		panic(err)
	}
	var myIp IP
	err = json.Unmarshal(body, &myIp)
	if err != nil {
		panic(err)
	}
	ipStr := "http://" + myIp.Query + ":8080"
	requestBody := strings.NewReader(ipStr)
	req, err := http.NewRequest(http.MethodGet, ip+"/addPeer", requestBody)
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		Log("Failed to connect to peer.", true)
	}
}
