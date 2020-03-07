package utils

// utils contains utility functions that are used in packages
import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/big"
	"os/user"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/crypto/sha3"
)

// ErrTypeNotSupported is an error returned if the given type conversion isn't supported yet
var ErrTypeNotSupported = errors.New("type not supported, please feel free to PR")

// Timestamp gets the human readable timestamp
func Timestamp() string {
	return time.Now().Format(time.RFC850)
}

// Unix gets the unix timestamp
func Unix() int64 {
	return time.Now().Unix()
}

// SHA3hash gets the SHA3-512 hash of the passed string
func SHA3hash(inputString string) string {
	byteString := sha3.Sum512([]byte(inputString))
	return hex.EncodeToString(byteString[:])
	// so now we have a SHA3hash that we can use to assign unique ids to our assets
}

// GetHomeDir gets the home directory of the user
func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		log.Println("error while getting current directory")
	}
	return usr.HomeDir, err
}

// GetRandomString gets a random string of length n
func GetRandomString(n int) string {
	// random string implementation courtesy: icza
	// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var err error
	int64Lim := new(big.Int).SetInt64(int64(1<<63 - 1)) // int64 lim = 2*63 -1
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, remain := n-1, letterIdxMax; i >= 0; {
		var tmp *big.Int
		tmp, err = rand.Int(rand.Reader, int64Lim)
		cache := tmp.Int64()
		if remain == 0 {
			tmp, err = rand.Int(rand.Reader, int64Lim)
			cache = tmp.Int64()
			remain = letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	if err != nil {
		b = []byte{0}
	}

	return string(b)
}

// ToBigInt converts a passed interface to big.Int
func ToBigInt(x interface{}) (*big.Int, error) {
	log.Println("calling anything to byte function")
	switch x.(type) {
	case big.Int:
		return x.(*big.Int), nil
	case []byte:
		return new(big.Int).SetBytes(x.([]byte)), nil
	case int:
		return big.NewInt(int64(x.(int))), nil
	case uint64:
		return new(big.Int).SetUint64(x.(uint64)), nil
	case string:
		log.Println([]byte(x.(string)))
		return new(big.Int).SetBytes([]byte(x.(string))), nil
	default:
		return new(big.Int).SetUint64(0), nil
	}
}

// ToByte converts a passed interface to bytes
func ToByte(x interface{}) ([]byte, error) {
	switch x.(type) {
	case []byte:
		return x.([]byte), nil
	case int:
		string1 := strconv.Itoa(x.(int))
		return []byte(string1), nil
	case uint32:
		a := make([]byte, 4)
		binary.BigEndian.PutUint32(a, x.(uint32))
		return a, nil
	case uint16:
		a := make([]byte, 2)
		binary.BigEndian.PutUint16(a, x.(uint16))
		return a[1:], nil
	case uint64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, x.(uint64))
		return b, nil
	}

	log.Println("type conversion not supported")
	return nil, ErrTypeNotSupported
}

// ToString converts a passed interface to string
func ToString(x interface{}) (string, error) {
	switch x.(type) {
	case string:
		return x.(string), nil
	case float64:
		return strconv.FormatFloat(x.(float64), 'f', 6, 64), nil
	case int64:
		return strconv.FormatInt(x.(int64), 10), nil // s == "97" (decimal)
	case int:
		return strconv.Itoa(x.(int)), nil
	}

	log.Println("type conversion not supported")
	return "", ErrTypeNotSupported
}

// ToInt converts a passed interface to int
func ToInt(x interface{}) (int, error) {
	switch x.(type) {
	case int:
		return x.(int), nil
	case string:
		return strconv.Atoi(x.(string))
	case []byte:
		return strconv.Atoi(string(x.([]byte)))
	}

	log.Println("type conversion not supported")
	return -1, ErrTypeNotSupported
}

// ToFloat converts a passed interface to float
func ToFloat(x interface{}) (float64, error) {
	switch x.(type) {
	case float64:
		return x.(float64), nil
	case []byte:
		return ToFloat(string(x.([]byte)))
	case string:
		return strconv.ParseFloat(x.(string), 32)
	case int:
		return float64(x.(int)), nil
	}

	log.Println("type conversion not supported")
	return -1, ErrTypeNotSupported
}

// ToUint16 converts a passed interface to uint16
func ToUint16(x interface{}) (uint16, error) {
	switch x.(type) {
	case uint16:
		return x.(uint16), nil
	case []byte:
		b := x.([]byte)
		if len(b) == 1 {
			zero := make([]byte, 1)
			b = append(zero, b...)
		}
		return binary.BigEndian.Uint16(b), nil
	}

	log.Println("type conversion not supported")
	return 0, ErrTypeNotSupported
}
