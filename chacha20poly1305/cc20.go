package cc20

import (
	"github.com/pkg/errors"
	"log"

	utils "github.com/Varunram/essentials/utils"
	cc20 "golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(input []byte, passphrase string) ([]byte, error) {
	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])
	aead, err := cc20.NewX(key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to instantiate XChaCha20-Poly1305")
	}

	nonce := make([]byte, cc20.NonceSizeX)
	nonce = []byte(utils.SHA3hash(sha3Hash))[0:24]

	ciphertext := aead.Seal(nonce, nonce, input, nil) // prepend nonce to the ciphertext
	return ciphertext, nil
}

func Decrypt(input []byte, passphrase string) ([]byte, error) {
	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])

	aead, err := cc20.NewX(key)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to instantiate XChaCha20-Poly1305")
	}

	nonce := make([]byte, cc20.NonceSizeX)
	nonce = []byte(utils.SHA3hash(sha3Hash))[0:24]

	copy(nonce, []byte(utils.SHA3hash(passphrase)[73:96]))

	nonce, ciphertext := input[:cc20.NonceSizeX], input[cc20.NonceSizeX:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal("Failed to decrypt or authenticate message:", err)
	}

	return plaintext, nil
}
