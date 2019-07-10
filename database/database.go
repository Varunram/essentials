package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"

	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)

func CreateDirs(dirs ...string) {
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// directory does not exist, create one
			os.MkdirAll(dir, os.ModePerm)
		}
	}
}

// don't lock since boltdb can only process one operation at a time. As the application
// grows bigger, this would be a major reason to search for a new db system

// OpenDB opens the db
func CreateDB(dir string, buckets ...[]byte) (*bolt.DB, error) {
	// we need to check and create this directory if it doesn't exist
	db, err := bolt.Open(dir, 0600, nil) // store this in its ownd database
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	for _, bucket := range buckets {
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(bucket) // the projects bucket contains all our projects
			if err != nil {
				log.Println("Error while creating projects bucket", err)
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

func OpenDB(dir string) (*bolt.DB, error) {
	return bolt.Open(dir, 0600, nil) // store this in its ownd database
}

// DeleteKeyFromBucket deletes a given key from the bucket bucketName but doesn
// not shift indices of elements succeeding the deleted element's index
func DeleteKeyFromBucket(dir string, key int, bucketName []byte) error {
	// deleting project might be dangerous since that would mess with the other
	// functions, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB(dir)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		b.Delete(utils.ItoB(key))
		return nil
	})
}

// Save inserts a passed Investor object into the database
func Save(dir string, bucketName []byte, x interface{}, key int) error {
	db, err := OpenDB(dir)
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		encoded, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "error while marshaling json struct")
		}
		return b.Put([]byte(utils.ItoB(key)), encoded)
	})
	return err
}

func Retrieve(dir string, bucketName []byte, key int) (interface{}, error) {
	var x interface{}
	db, err := OpenDB(dir)
	if err != nil {
		return x, errors.Wrap(err, "failed to open db")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		byteString := b.Get(utils.ItoB(key))
		if byteString == nil {
			// no investor with the specific details
			return errors.New("No investor found with required credentials")
		}
		return json.Unmarshal(byteString, x)
	})
	return x, err
}

func RetrieveAllKeys(dir string, bucketName []byte) ([][]byte, error) {
	var arr [][]byte
	db, err := OpenDB(dir)
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName) // the projects bucket contains all our projects
		if err != nil {
			log.Println("Error while creating projects bucket", err)
			return err
		}
		return nil
	})

	if err != nil {
		return arr, errors.New("bucket not created and error while creating new bucket")
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for i := 1; ; i++ {
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			arr = append(arr, x)
		}
	})
	return arr, err
}
