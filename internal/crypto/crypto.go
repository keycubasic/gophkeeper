// Package crypto contains cryptographic functions.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// Encode encrypts and encodes message.
func Encode(secret, msg string) (string, error) {
	aes, err := aes.NewCipher([]byte(secret)[:32])
	if err != nil {
		return "", fmt.Errorf("crypto encode cypher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("crypto encode gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("crypto encode rand: %w", err)
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(msg), nil)
	encoded := hex.EncodeToString(encrypted)

	return encoded, nil
}

// Decode dencodes and decrypts message.
func Decode(secret, encoded string) (string, error) {
	encrypted, err := hex.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("crypto decode hex: %w", err)
	}

	aes, err := aes.NewCipher([]byte(secret)[:32])
	if err != nil {
		return "", fmt.Errorf("crypto decode cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", fmt.Errorf("crypto decode gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	decrypted, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", fmt.Errorf("crypto decode open: %w", err)
	}

	return string(decrypted), nil
}
