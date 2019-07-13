package scan

// package scan is not in utils since we can't test the below functions (which require
// user interaction) whereas functions in utils are essential
// and need to be tested in order for stuff to run properly

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"syscall"

	utils "github.com/Varunram/essentials/utils"
	"golang.org/x/crypto/ssh/terminal"
)

// ScanForInt scans for an integer
func ScanForInt() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, errors.New("Couldn't read user input")
	}
	num := scanner.Text()
	numI, err := strconv.Atoi(num)
	if err != nil {
		return -1, errors.New("Input not a number")
	}
	return numI, nil
}

// ScanForFloat scans for a float
func ScanForFloat() (float64, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, errors.New("Couldn't read user input")
	}
	num := scanner.Text()
	x, err := strconv.ParseFloat(num, 32)
	// ignore this error since we hopefully call this in the right place
	return x, err
}

// ScanForString scans for a string
func ScanForString() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	return inputString, nil
}

// ScanForStringWithCheckI scans for a string checking whether it is an integer
func ScanForStringWithCheckI() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	_, err := strconv.Atoi(inputString) // check whether input string is a number (for payback)
	if err != nil {
		return "", err
	}
	return inputString, nil
}

// ScanForStringWithCheckF scans for a string checking whether its a float
func ScanForStringWithCheckF() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	_, err := utils.ToFloat(inputString)
	if err != nil {
		fmt.Println("Amount entered is not a float, quitting")
		return "", errors.New("Amount entered is not a float, quitting")
	}
	return inputString, nil
}

// ScanForPassword scans for a password
func ScanForPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	tempString := string(bytePassword)
	hashedPassword := utils.SHA3hash(tempString)
	return hashedPassword, nil
}

// ScanRawPassword scans for a raw password
func ScanRawPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}
