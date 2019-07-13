// +build all travis

package utils

import (
	"log"
	"os/user"
	"testing"
	"time"
)

func TestUtils(t *testing.T) {
	// test out stuff here
	testString := "10"
	testStringFloat := "10.000000"
	testInt := 10
	testByte := []byte("10")
	testFloat := 10.0
	var err error

	testSlice, err := ToByte(10)
	if err != nil {
		t.Fatal(err)
	}

	// slcies can be compared only with nil
	for i, char := range testSlice {
		if char != testByte[i] {
			t.Fatalf("ItoB deosn't work as expected, quitting!")
		}
	}

	testSlice = []byte("01")
	check := false
	for i, char := range testSlice {
		if char == testByte[i] {
			check = true
		}
	}

	if check {
		t.Fatalf("Failed to catch error while comparing two different slcies")
	}

	x1, err := ToString(testInt)
	if err != nil || x1 != testString {
		t.Fatalf("ToString deosn't work as expected, quitting!")
	}

	x2, err := ToString(testInt)
	if err != nil || x2 == "" {
		t.Fatalf("ToString deosn't work as expected, quitting!")
	}

	x3, err := ToInt(testByte)
	if err != nil || x3 != testInt {
		log.Println(x3, err)
		t.Fatalf("ToInt deosn't work as expected, quitting!")
	}

	x4, err := ToInt(testByte)
	if err == nil && x4 == 9 {
		t.Fatalf("ToInt deosn't work as expected, quitting!")
	}

	x5, err := ToString(testFloat)
	if err != nil || x5 != testStringFloat {
		t.Fatalf("ToString deosn't work as expected, quitting!")
	}

	x6, err := ToString(testFloat)
	if err == nil && x6 == testString {
		t.Fatalf("ToString deosn't work as expected, quitting!")
	}

	x7, err := ToFloat(testStringFloat)
	if err != nil || x7 != testFloat {
		t.Fatalf("StoF deosn't work as expected, quitting!")
	}

	x8, err := ToFloat(testStringFloat)
	if err == nil && x8 == 9.0 {
		t.Fatalf("StoF deosn't work as expected, quitting!")
	}

	x9, err := ToInt(testString)
	if err != nil || x9 != testInt {
		t.Fatalf("StoI deosn't work as expected, quitting!")
	}

	x10, err := ToInt(testString)
	if err == nil && x10 == 9 {
		t.Fatalf("StoI deosn't work as expected, quitting!")
	}

	if SHA3hash("password") != "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716" {
		t.Fatalf("SHA3 doesn't work as expected, quitting!")
	}

	if SHA3hash("blah") == "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716" {
		t.Fatalf("SHA3 doesn't work as expected, quitting!")
	}

	hd, err := GetHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	usr, err := user.Current()
	if usr.HomeDir != hd || err != nil {
		t.Fatalf("Home directories don't match, quitting!")
	}

	rs := GetRandomString(10)
	if len(rs) != 10 {
		t.Fatalf("Random string length not equal to what is expected")
	}

	if time.Now().Format(time.RFC850) != Timestamp() {
		t.Fatalf("Timestamps don't match, quitting!")
	}

	if time.Now().Unix() != Unix() {
		t.Fatalf("Timestamps don't match, quitting!")
	}

	y, err := ToString(123412341234)
	if err != nil || y != "123412341234" {
		t.Fatalf("I64 to string doesn't work, quitting!")
	}

	_, err = ToFloat("blah")
	if err == nil {
		t.Fatalf("Not able to catch invalid string to float error!")
	}

	_, err = ToInt("blah")
	if err == nil {
		t.Fatalf("Not able to catch invalid string to float error!")
	}
}
