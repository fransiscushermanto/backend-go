package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Encrypt uses AES-GCM to encrypt plaintext with a given key.
// It returns a single byte slice containing the nonce and the ciphertext.
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	// Create a new AES cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher mode instance.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Create the nonce. GCM's standard nonce size is 12 bytes.
	nonce := make([]byte, gcm.NonceSize())
	// Populate the nonce with cryptographically secure random data.
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt the data. Seal will append the ciphertext to the nonce and return
	// the combined slice. The first argument is the destination slice, which
	// we set to the nonce itself to achieve nonce + ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func Decrypt(key []byte, dataToDecrypt []byte) ([]byte, error) {
	// Create a new AES cipher block from the key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create a GCM cipher mode instance.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// The standard nonce size for GCM is 12 bytes.
	nonceSize := gcm.NonceSize()
	if len(dataToDecrypt) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Split the data into the nonce and the actual encrypted message.
	nonce, ciphertext := dataToDecrypt[:nonceSize], dataToDecrypt[nonceSize:]

	// Decrypt the data. If the key is wrong or the data is corrupt, it will return an error.
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
