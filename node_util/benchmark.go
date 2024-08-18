package node_util

import (
	"fmt"
	"time"
)

// Benchmark is a function that performs a benchmark by creating transactions, creating a block, and measuring the time it takes.
//
// It creates two transactions using the GetKey function to get the public key of the sender and recipient. The transactions have a zero amount, a zero sender signature, a timestamp of the current time, no contracts, no body, no body signatures, and a flag indicating that it is not from a smart contract. The transactions are appended to the MiningTransactions slice.
//
// It then starts a timer using the time.Now function.
//
// It calls the CreateBlock function to create a block. If an error occurs during the creation of the block, it panics.
//
// It stops the timer using the time.Since function and calculates the elapsed time.
//
// It calculates the score by dividing the block's difficulty by the elapsed time in seconds.
//
// It prints the elapsed time and the score using the fmt.Printf function.
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
