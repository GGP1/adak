package crypt

import (
	"bytes"
	"testing"

	"github.com/spf13/viper"
)

func TestCrypt(t *testing.T) {
	viper.Set("token.secretKey", "Aj _'X0#Zea8w@2")
	data := []byte("testing crypt package")

	ciphertext, err := Encrypt(data)
	if err != nil {
		t.Fatalf("Failed encrypting data: %v", err)
	}

	plaintext, err := Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Failed decrypting data: %v", err)
	}

	if !bytes.Equal(plaintext, data) {
		t.Errorf("Expected %q, got %q", data, plaintext)
	}
}
