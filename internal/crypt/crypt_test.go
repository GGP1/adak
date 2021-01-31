package crypt

import (
	"bytes"
	"testing"
)

func TestCrypt(t *testing.T) {
	key := []byte("Aj _'X0#Zea8w@2")
	data := []byte("testing crypt package")

	ciphertext, err := Encrypt(key, data)
	if err != nil {
		t.Fatalf("Failed encrypting data: %v", err)
	}

	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Failed decrypting data: %v", err)
	}

	if !bytes.Equal(plaintext, data) {
		t.Errorf("Expected %q, got %q", data, plaintext)
	}
}
