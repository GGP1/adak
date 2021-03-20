package crypt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	// Do not provide additional information about the failure to potential attackers
	errEncrypt = errors.New("encrypt error")
	errDecrypt = errors.New("decrypt error")
)

// Encrypt ciphers data with the given key.
func Encrypt(data []byte) ([]byte, error) {
	hash := createHMAC()

	AEAD, err := chacha20poly1305.New(hash)
	if err != nil {
		return nil, errEncrypt
	}

	nonce := make([]byte, AEAD.NonceSize())

	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errEncrypt
	}

	dst := make([]byte, AEAD.NonceSize())
	copy(dst, nonce)

	ciphertext := AEAD.Seal(dst, nonce, data, nil)

	return ciphertext, nil
}

// Decrypt deciphers data with the given key.
func Decrypt(data []byte) ([]byte, error) {
	hash := createHMAC()

	AEAD, err := chacha20poly1305.New(hash)
	if err != nil {
		return nil, errDecrypt
	}

	nonceSize := AEAD.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := AEAD.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errDecrypt
	}

	return plaintext, nil
}

// Create an HMAC SHA256 hash (32 bytes) with the key provided.
func createHMAC() []byte {
	key := []byte(viper.GetString("token.secretkey"))
	hash := hmac.New(sha256.New, key)

	return hash.Sum(nil)
}
