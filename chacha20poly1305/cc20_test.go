package cc20

import (
	"log"
	"testing"
)

func TestCC20(t *testing.T) {
	password := "Cool"
	ciphertext, err := Encrypt([]byte("Hello World"), password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Encrypted: %x\n", ciphertext)
	plaintext, err := Decrypt(ciphertext, password)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Decrypted: %s\n", plaintext)
}
