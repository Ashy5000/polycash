package main

var blockchain []Block

func Append(block Block) {
	blockchain = append(blockchain, block)
}
