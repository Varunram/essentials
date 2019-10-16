package cc20

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	cc20 "golang.org/x/crypto/chacha20poly1305"
)

// package cc20 implements encrypt/decrypt functions for the chacha20poly1305 cipher
// we implement XChaCha20 which is a variant of CC20 allowing random nonces to be safe
// https://libsodium.gitbook.io/doc/advanced/stream_ciphers/xchacha20

// Encrypt encrypts a given passphrase using CC20-poly1305
func Encrypt(input []byte, passphrase string) ([]byte, error) {
	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])
	aead, err := cc20.NewX(key)
	if err != nil {
		log.Println("Failed to instantiate XChaCha20-Poly1305", err)
		return nil, errors.Wrap(err, "Failed to instantiate XChaCha20-Poly1305")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:24]

	ciphertext := aead.Seal(nil, nonce, input, nil) // prepend nonce to the ciphertext
	return ciphertext, nil
}

// Decrypt decrypts a given cipher with the passed passphrase
func Decrypt(input []byte, passphrase string) ([]byte, error) {
	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])

	aead, err := cc20.NewX(key)
	if err != nil {
		log.Println("Failed to instantiate XChaCha20-Poly1305", err)
		return nil, errors.Wrap(err, "Failed to instantiate XChaCha20-Poly1305")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:24]

	plaintext, err := aead.Open(nil, nonce, input, nil)
	if err != nil {
		log.Println("failed to decrypt or authenticate message: ", err)
		return nil, errors.Wrap(err, "failed to decrypt or authenticate  message")
	}

	return plaintext, nil
}
