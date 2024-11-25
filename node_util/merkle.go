package node_util

import (
	"crypto/sha256"
	"encoding/hex"
)

type MerkleNode struct {
	Data     []byte
	Children map[byte]int
	Parent   int
	Hash     string
	Key      string
}

func HashNode(tree []MerkleNode, nodeIndex int) string {
	hasher := sha256.New()
	hasher.Write(tree[nodeIndex].Data)
	for _, index := range tree[nodeIndex].Children {
		hasher.Write([]byte(tree[index].Hash))
	}
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func HashTree(tree []MerkleNode, nodeIndex int) []MerkleNode {
	// Base case
	if len(tree[nodeIndex].Children) == 0 {
		tree[nodeIndex].Hash = HashNode(tree, nodeIndex)
	}
	// Recursive case
	for _, index := range tree[nodeIndex].Children {
		tree = HashTree(tree, index)
	}
	tree[nodeIndex].Hash = HashNode(tree, nodeIndex)
	return tree
}

func InsertValue(tree []MerkleNode, key string, value []byte) []MerkleNode {
	activeIndex := 0
	activeKey := ""
	if len(tree) == 0 {
		tree = append(tree, MerkleNode{
			Data:     make([]byte, 0),
			Children: make(map[byte]int),
			Parent:   0,
			Hash:     "",
			Key:      "",
		})
		activeIndex = 0
	}
	for i, _ := range key {
		activeKey += string(key[i])
		if _, ok := tree[activeIndex].Children[key[i]]; ok {
			activeIndex = tree[activeIndex].Children[key[i]]
			continue
		}
		tree[activeIndex].Children[key[i]] = len(tree)
		tree = append(tree, MerkleNode{
			Data:     make([]byte, 0),
			Children: make(map[byte]int),
			Parent:   activeIndex,
			Hash:     "",
			Key:      activeKey,
		})
		activeIndex = len(tree) - 1
	}
	tree[activeIndex].Data = []byte(value)
	tree = HashTree(tree, activeIndex)
	return tree
}

func GetValue(tree []MerkleNode, key string) ([]byte, bool) {
	activeIndex := 0
	for i, _ := range key {
		if val, ok := tree[activeIndex].Children[key[i]]; ok {
			activeIndex = val
		} else {
			return []byte(""), false
		}
	}
	return tree[activeIndex].Data, true
}

func Merge(initial []MerkleNode, delta []MerkleNode) []MerkleNode {
	// Get values from delta
	for _, node := range delta {
		if node.Data != nil {
			initial = InsertValue(initial, node.Key, node.Data)
		}
	}
	return initial
}
