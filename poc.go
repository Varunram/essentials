package main

import (
	"encoding/json"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

// Golang segfault bug
var bucketName = []byte("test")
var dir = "test.db"

type Test struct {
	Name string
}

func ItoB(a int) []byte {
	string1 := strconv.Itoa(a)
	return []byte(string1)
}

func populateJunkValues(x Test, key int) error {
	db, err := bolt.Open(dir, 0600, nil)
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
		return b.Put([]byte(ItoB(key)), encoded)
	})
	return err
}

func test1() {
	arr := make(map[int][]byte)
	var x Test
	x.Name = "name"
	for i := 1; i < 40; i++ {
		err := populateJunkValues(x, i)
		if err != nil {
			log.Fatal(err)
		}
	}
	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		log.Fatal("Error while opening database")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			log.Fatal("Error while creating projects bucket")
		}
		return nil
	})

	if err != nil {
		log.Fatal("error while creating db, exiting")
	}

	db.Close()

	db, err = bolt.Open(dir, 0600, nil)
	if err != nil {
		log.Fatal("Error while opening database")
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for i := 1; ; i++ {
			x := b.Get(ItoB(i))
			if x == nil {
				return nil
			}
			arr[i] = x
			log.Println("THIS WORKS FINE: ", arr[i])
		}
	})

	db.Close()
	log.Println("SEGFAULT:")
	for i := 1; i < len(arr); i++ {
		log.Println(string(arr[i]))
	}

	log.Println(arr)
}

func test2() {
	arr := make(map[int][]byte)
	var x Test
	x.Name = "name"
	for i := 1; i < 40; i++ {
		err := populateJunkValues(x, i)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		log.Fatal("Error while opening database")
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for i := 1; ; i++ {
			x := b.Get(ItoB(i))
			if x == nil {
				return nil
			}
			arr[i] = x
			log.Println("THIS WORKS FINE: ", arr[i])
		}
	})

	db.Close()
	log.Println("SEGFAULT:")
	for i := 1; i < len(arr); i++ {
		log.Println(string(arr[i]))
	}

	log.Println(arr)
}

func test3() []byte {
	var arr []byte
	var x Test
	x.Name = "name"
	for i := 1; i < 40; i++ {
		err := populateJunkValues(x, i)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err := bolt.Open(dir, 0600, nil)
	if err != nil {
		log.Fatal("Error while opening database")
	}

	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for i := 1; ; i++ {
			x := b.Get(ItoB(i))
			if x == nil {
				return nil
			}
			arr = x
			log.Println("THIS WORKS FINE: ", arr)
		}
	})

	log.Println("SEGFAULT:")
	return arr
}

func main() {
	// test1()
	// test2()
	arr := test3()
	log.Println(arr)
}
