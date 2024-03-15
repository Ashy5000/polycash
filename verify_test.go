package main

import (
	"crypto/dsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strconv"
	"testing"
	"time"
)

func TestVerifyTransaction(t *testing.T) {
	t.Run("It should return true if the transaction is valid", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		sender := key.Y
		receiver := key.Y
		amount := "0"
		r := big.NewInt(0)
		s := big.NewInt(0)
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender, receiver, amount)))
		r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
		if err != nil {
			panic(err)
		}
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, *r, *s)
		assert.True(t, result)
	})
	t.Run("It should return false if the transaction double spends", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		sender := key.Y
		receiver := key.Y
		amount := "1"
		r := big.NewInt(0)
		s := big.NewInt(0)
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender, receiver, amount)))
		r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
		if err != nil {
			panic(err)
		}
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, *r, *s)
		assert.False(t, result)
	})
	t.Run("It should return false if the transaction signature is invalid", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		amount := "1"
		r := big.NewInt(0)
		s := big.NewInt(0)
		result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, *r, *s)
		assert.False(t, result)
	})
}

func TestVerifyMiner(t *testing.T) {
	t.Run("It should return true if the miner is valid", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		result := VerifyMiner(key.PublicKey)
		assert.True(t, result)
	})
	t.Run("It should return false if there are too many miners", func(t *testing.T) {
		key := GetKey()
		miner := key.PublicKey
		miner.Y = big.NewInt(1)
		blockchain = nil
		Append(GenesisBlock())
		Append(Block{
			Sender:            key.PublicKey,
			Recipient:         key.PublicKey,
			Miner:             miner,
			Amount:            0,
			PreviousBlockHash: HashBlock(GenesisBlock()),
			Difficulty:        1,
		})
		result := VerifyMiner(key.PublicKey)
		assert.False(t, result)
	})
}

func TestVerifyBlock(t *testing.T) {
	t.Run("It should return true if the block is valid", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey
		receiver := key.PublicKey
		amount := 0.0
		r := big.NewInt(0)
		s := big.NewInt(0)
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender.Y, receiver.Y, strconv.FormatFloat(amount, 'f', -1, 64))))
		r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
		if err != nil {
			panic(err)
		}
		block, err := CreateBlock(sender, receiver, amount, *r, *s, hash, strconv.FormatInt(time.Now().UnixNano(), 10))
		if err != nil {
			panic(err)
		}
		result := VerifyBlock(block)
		assert.True(t, result)
	})
	t.Run("It should return false if the block has an invalid transaction", func(t *testing.T) {
		key := GetKey()
		blockchain = nil
		Append(GenesisBlock())
		sender := key.PublicKey
		receiver := key.PublicKey
		amount := "1"
		r := big.NewInt(0)
		s := big.NewInt(0)
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s", sender.Y, receiver.Y, amount)))
		r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
		if err != nil {
			panic(err)
		}
		block := Block{
			Sender:            sender,
			Recipient:         receiver,
			Amount:            0,
			R:                 *r,
			S:                 *s,
			PreviousBlockHash: HashBlock(GenesisBlock()),
			Miner:             key.PublicKey,
			Difficulty:        1,
		}
		result := VerifyBlock(block)
		assert.False(t, result)
	})
}
