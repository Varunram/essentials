package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"runtime"

	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)

var ErrBucketMissing = errors.New("bucket doesn't exist")
var ErrElementNotFound = errors.New("element not found")

// package database contains useful boltdb handlers

// CreateDirs creates (db) directories if they don't exist
func CreateDirs(dirs ...string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// directory does not exist, create one
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}

// CreateDB creates a new database
func CreateDB(dir string, buckets ...[]byte) (*bolt.DB, error) {
	// we need to check and create this directory if it doesn't exist

	// OpenDB creates a db if it doens't exist yet
	db, err := OpenDB(dir)
	if err != nil {
		return db, errors.New("Couldn't open database, exiting!")
	}
	for _, bucket := range buckets {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return db, errors.Wrap(err, "could not create bucket")
		}
	}
	return db, nil
}

// OpenDB opens the database
func OpenDB(dir string) (*bolt.DB, error) {
	if runtime.GOOS == "linux" {
		return bolt.Open(dir, 0600, &bolt.Options{
			MmapFlags: 0x8000, // MAP_POPULATE = 0x8000
		})
	}
	return bolt.Open(dir, 0600, nil)
}

// DeleteKeyFromBucket deletes a given key from a bucket
func DeleteKeyFromBucket(dir string, key int, bucketName []byte) error {
	db, err := OpenDB(dir)
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	// if the passed key is not integer, don't open the db
	iK, err := utils.ToByte(key)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}
		b.Delete(iK)
		return nil
	})
}

// Save inserts an interface with an integer key
func Save(dir string, bucketName []byte, x interface{}, key int) error {
	db, err := OpenDB(dir)
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	// if x is not interace, don't open the database
	encoded, err := json.Marshal(x)
	if err != nil {
		return errors.Wrap(err, "error while marshaling json struct")
	}

	// if the passed key is not integer, don't open the db
	iK, err := utils.ToByte(key)
	if err != nil {
		return err
	}

	// open the db only to insert the element and don't check for other stuff
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}
		return b.Put(iK, encoded)
	})

	return err
}

// Retrieve retrieves a byteStringf rom the database
func Retrieve(dir string, bucketName []byte, key int) ([]byte, error) {
	var returnBytes []byte
	db, err := OpenDB(dir)
	if err != nil {
		return returnBytes, errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	// if the passed key is not integer, don't open the db
	iK, err := utils.ToByte(key)
	if err != nil {
		return returnBytes, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}
		x := b.Get(iK)
		if x == nil {
			return ErrElementNotFound
		}
		returnBytes = make([]byte, len(x))
		copy(returnBytes, x)
		return nil
	})

	return returnBytes, err
}

// RetrieveAllKeys retrieves all key value pairs from the database
func RetrieveAllKeys(dir string, bucketName []byte) ([][]byte, error) {
	var arr [][]byte
	db, err := OpenDB(dir)
	if err != nil {
		return arr, errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}
		for i := 1; ; i++ {
			iB, err := utils.ToByte(i)
			if err != nil {
				return err
			}
			x := b.Get(iB)
			if x == nil {
				return nil
			}
			temp := make([]byte, len(x))
			copy(temp, x)
			arr = append(arr, temp)
		}
	})
	return arr, err
}
