// Copyright 2024, Asher Wrobel
/*
This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/dsa"
    "crypto/rand"
    "encoding/json"
    "io"
    "os"
)

func IsKeyEncrypted() bool {
    // Check if key.json is encrypted.
    contents, err := os.ReadFile("key.json")
    if err != nil {
        Error("No key found.", true)
    }
    var key dsa.PrivateKey
    err = json.Unmarshal(contents, &key)
    if err != nil {
        return true
    }
    return false
}

func EncryptKey(password string) {
    plaintext, err := os.ReadFile("key.json")
    if err != nil {
        Error("No key found.", true)
    }
    block, err := aes.NewCipher([]byte(password))
    if err != nil {
        Error("Error creating cipher. Ensure that the password is a multiple of 16 characters long.", false)
        return
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        panic(err)
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        panic(err)
    }
    cipherText := gcm.Seal(nonce, nonce, plaintext, nil)
    err = os.WriteFile("key.json", cipherText, 0644)
    if err != nil {
        panic(err)
    }
}

func DecryptKey(password string) {
    ciphertext, err := os.ReadFile("key.json")
    if err != nil {
        panic(err)
    }
    block, err := aes.NewCipher([]byte(password))
    if err != nil {
        panic(err)
    }
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        panic(err)
    }
    nonceSize := gcm.NonceSize()
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        panic(err)
    }
    err = os.WriteFile("key.json", plaintext, 0644)
    if err != nil {
        panic(err)
    }
}
