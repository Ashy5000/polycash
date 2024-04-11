// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
    "crypto/dsa"
    "crypto/rand"
    "crypto/sha256"
    "fmt"
    "math/big"
    "strconv"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
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
        timestamp := time.Now()
        hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", sender, receiver, amount, timestamp.UnixNano())))
        r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
        if err != nil {
            panic(err)
        }
        result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, timestamp, *r, *s)
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
        result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, time.Now(), *r, *s)
        assert.False(t, result)
    })
    t.Run("It should return false if the transaction signature is invalid", func(t *testing.T) {
        key := GetKey()
        blockchain = nil
        Append(GenesisBlock())
        amount := "1"
        r := big.NewInt(0)
        s := big.NewInt(0)
        result := VerifyTransaction(key.PublicKey, key.PublicKey, amount, time.Now(), *r, *s)
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
            Transactions:      []Transaction{},
            Miner:             miner,
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
        timestamp := time.Now()
        hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%d", sender.Y, receiver.Y, strconv.FormatFloat(amount, 'f', -1, 64), timestamp.UnixNano())))
        r, s, err := dsa.Sign(rand.Reader, &key, hash[:])
        if err != nil {
            panic(err)
        }

        miningTransactions = []Transaction{
            {
                Sender:    sender,
                Recipient: receiver,
                Amount:    amount,
                SenderSignature: Signature{
                    R: *r,
                    S: *s,
                },
        Timestamp: timestamp,
            },
        }
        block, err := CreateBlock()
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
            Transactions: []Transaction{
                {
                    Sender:          sender,
                    Recipient:       receiver,
                    Amount:          0,
                    SenderSignature: Signature{R: *r, S: *s},
                },
            },
            PreviousBlockHash: HashBlock(GenesisBlock()),
            Miner:             key.PublicKey,
            Difficulty:        1,
        }
        result := VerifyBlock(block)
        assert.False(t, result)
    })
}
