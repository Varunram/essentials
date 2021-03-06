package database

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)

// ErrBucketMissing is an error handler for a missing bucket
var ErrBucketMissing = errors.New("bucket doesn't exist")

// ErrElementNotFound is an error handler for a missing element
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
		log.Println("couldn't open database, exiting: ", err)
		return db, errors.New("couldn't open database, exiting")
	}

	for _, bucket := range buckets {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket)
			return err
		})
		if err != nil {
			log.Println("could not create bucket: ", err)
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

	defer func() {
		if ferr := db.Close(); ferr != nil {
			err = ferr
		}
	}()

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

	defer func() {
		if ferr := db.Close(); ferr != nil {
			err = ferr
		}
	}()

	// if x is not interace, don't open the database
	encoded, err := json.Marshal(x)
	if err != nil {
		log.Println("error while marshaling json struct: ", err)
		return errors.Wrap(err, "error while marshaling json struct")
	}

	// if the passed key is not integer, don't open the db
	iK, err := utils.ToByte(key)
	if err != nil {
		return err
	}

	// open the db only to insert the element and don't check for other stuff
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return b.Put(iK, encoded)
	})

	return err
}

// Retrieve retrieves a byteString from the database
func Retrieve(dir string, bucketName []byte, key int) ([]byte, error) {
	var returnBytes []byte
	db, err := OpenDB(dir)
	if err != nil {
		return returnBytes, errors.Wrap(err, "could not open database")
	}

	defer func() {
		if ferr := db.Close(); ferr != nil {
			err = ferr
		}
	}()

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

	defer func() {
		if ferr := db.Close(); ferr != nil {
			err = ferr
		}
	}()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}
		err := b.ForEach(func(k, x []byte) error {
			temp := make([]byte, len(x))
			copy(temp, x)
			arr = append(arr, temp)
			return nil
		})
		return err
	})
	return arr, err
}

// RetrieveAllKeysLim gets the total number of keys in a bucket
func RetrieveAllKeysLim(dir string, bucketName []byte) (int, error) {
	lim := 0
	db, err := OpenDB(dir)
	if err != nil {
		return lim, errors.Wrap(err, "could not open database")
	}

	defer func() {
		if ferr := db.Close(); ferr != nil {
			err = ferr
		}
	}()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return ErrBucketMissing
		}

		x := b.Stats().KeyN
		var temp int
		temp = x // weird pointer stuff thanks to boltdb

		lim = temp
		return nil
	})

	if err != nil {
		log.Println("could not open db for reading: ", err)
	}

	return lim, err
}
