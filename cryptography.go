package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

func Unpad(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

// CBC
func encryptCBC(key, plaintext []byte) (ciphertext []byte, err error) {
	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	//iv, _ := hex.DecodeString("acfa7a047800b2f221f2c4f7d626eafb")
	//copy(ciphertext[:aes.BlockSize], iv)

	fmt.Printf("CBC Key: %s\n", hex.EncodeToString(key))
	fmt.Printf("CBC IV: %s\n", hex.EncodeToString(iv))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}
func decryptCBC(key, ciphertext []byte) (plaintext []byte, err error) {
	var block cipher.Block

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		fmt.Printf("ciphertext too short")
		return
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	plaintext = ciphertext

	return
}

func aesmain() {
	var ciphertext, plaintext []byte
	var err error

	// The key length can be 32, 24, 16  bytes (OR in bits: 128, 192 or 256)
	key := []byte("longer means more possible keys ")
	plaintext = []byte("longer MEANS more POSSIBLE keys ")

	if ciphertext, err = encryptCBC(key, plaintext); err != nil {
		panic(err)
	}
	cleartext := base64.StdEncoding.EncodeToString(ciphertext[aes.BlockSize:])
	fmt.Printf("CBC: %s\n", cleartext)

	if plaintext, err = decryptCBC(key, ciphertext); err != nil {
		panic(err)
	}
	fmt.Printf("Clear from CBC: %s\n", plaintext)

	// echo 'dev' | openssl enc -aes-256-cbc -nosalt -K 6c6f6e676572206d65616e73206d6f726520706f737369626c65206b65797320 -iv acfa7a047800b2f221f2c4f7d626eafb | base64 + add IV padding first
	//enc, _ := base64.StdEncoding.DecodeString("rPp6BHgAsvIh8sT31ibq+xX5Kc/0uBTewwk+lCd5v0M=")
	enc, _ := base64.StdEncoding.DecodeString("Tz3w9bV7G/DOa2R8yCTU4LwmT/ahQ9sWRNEMtBXIUuQ=")
	test_key := []byte("HELP ME !!")
	// Pad key
	padded := make([]byte, 32)
	copy(padded, []byte(test_key))
	if plaintext, err = decryptCBC(padded, enc); err != nil {
		panic(err)
	}
	// We need to unpad plaintext
	fmt.Printf("padded key: %s plain: %s\n", padded, Unpad(plaintext))
}
