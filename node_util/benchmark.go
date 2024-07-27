package node_util

import (
	"fmt"
	"time"
)

func Benchmark() {
	// Create transactions
	for i := 0; i < 2; i++ {
		transaction := Transaction{
			Sender: PublicKey{
				Y: GetKey("").PublicKey.Y,
			},
			Recipient: PublicKey{
				Y: GetKey("").PublicKey.Y,
			},
			Amount:            0,
			SenderSignature:   Signature{},
			Timestamp:         time.Now(),
			Contracts:         nil,
			Body:              nil,
			BodySignatures:    nil,
			FromSmartContract: false,
		}
		MiningTransactions = append(MiningTransactions, transaction)
	}
	// Start timer
	start := time.Now()
	// Create a block
	block, err := CreateBlock()
	if err != nil {
		panic(err)
	}
	// Stop timer
	elapsed := time.Since(start)
	score := float64(block.Difficulty) / elapsed.Seconds()
	fmt.Printf("Time: %s\n", elapsed)
	fmt.Printf("Score: %f\n", score)
}
