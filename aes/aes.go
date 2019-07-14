package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"

	utils "github.com/Varunram/essentials/utils"
)

// the aes package implements AES-256 GCM encryption and decrpytion functions

// Encrypt encrypts a given data stream with a given passphrase
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])

	block, _ := aes.NewCipher(key)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.Wrap(err, "Error while opening new GCM block")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:gcm.NonceSize()]

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts a given data stream with a given passphrase
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 || len(passphrase) == 0 {
		return data, errors.New("Length of data is zero, can't decrpyt!")
	}

	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, errors.Wrap(err, "Error while initializing new cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return data, errors.Wrap(err, "failed to initialize new gcm block")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:gcm.NonceSize()]

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return plaintext, errors.Wrap(err, "failed to decrypt data")
	}

	return plaintext, nil
}

// EncryptFile encrypts a given file with the given passphrase
func EncryptFile(filename string, data []byte, passphrase string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, "Error while creating file")
	}
	defer f.Close()
	data, err = Encrypt(data, passphrase)
	if err != nil {
		return errors.Wrap(err, "Error while encrypting file")
	}
	f.Write(data)
	return nil
}

// DecryptFile encrypts a given file with the given passphrase
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}
