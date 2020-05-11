package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"

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
		log.Println("Error while opening new GCM block")
		return nil, errors.Wrap(err, "Error while opening new GCM block")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:gcm.NonceSize()]

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt decrypts a given data stream with a given passphrase
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	if len(data) == 0 || len(passphrase) == 0 {
		log.Println("length of data is zero, can't decrypt")
		return data, errors.New("length of data is zero, can't decrpyt")
	}

	sha3Hash := utils.SHA3hash(passphrase)
	key := []byte(sha3Hash[0:32])
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Error while initializing new cipher")
		return data, errors.Wrap(err, "Error while initializing new cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("failed to initialize new gcm block")
		return data, errors.Wrap(err, "failed to initialize new gcm block")
	}

	nonce := []byte(utils.SHA3hash(sha3Hash))[0:gcm.NonceSize()]

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		log.Println("failed to decrypt data")
		return plaintext, errors.Wrap(err, "failed to decrypt data")
	}

	return plaintext, nil
}

// EncryptFile encrypts a given file with the given passphrase
func EncryptFile(filename string, data []byte, passphrase string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Println("Error while creating file")
		return errors.Wrap(err, "Error while creating file")
	}

	defer func() {
		if ferr := f.Close(); ferr != nil {
			err = ferr
			os.Remove(filename)
		}
	}()

	data, err = Encrypt(data, passphrase)
	if err != nil {
		log.Println("Error while encrypting file")
		err = os.Remove(filename)
		if err != nil {
			log.Println("could not delete file: ", filename)
			panic(err) // panic since we can't return
		}
		return errors.Wrap(err, "Error while encrypting file")
	}

	f.Write(data)
	return nil
}

// DecryptFile encrypts a given file with the given passphrase
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("could not read from file: ", err)
		return data, err
	}
	dData, err := Decrypt(data, passphrase)
	if err != nil {
		log.Println("could not decrypt data: ", err)
	}
	return dData, err
}
