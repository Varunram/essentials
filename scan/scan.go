package scan

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"syscall"

	utils "github.com/Varunram/essentials/utils"
	"golang.org/x/crypto/ssh/terminal"
)

// package scan can be used by CLI clients that want to accept inptus from the CLI

// ScanInt scans for an integer
func ScanInt() (int, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return -1, errors.New("Couldn't read user input")
	}
	num := scanner.Text()
	numI, err := utils.ToInt(num)
	if err != nil {
		return -1, errors.New("Input not a number")
	}
	return numI, nil
}

// ScanFloat scans for a float
func ScanFloat() (float64, error) {
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

// ScanString scans for a string
func ScanString() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	return inputString, nil
}

// ScanStringCheckInt scans for a string checking whether it is an integer
func ScanStringCheckInt() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	_, err := utils.ToInt(inputString)
	if err != nil {
		return "", err
	}
	return inputString, nil
}

// ScanStringCheckFloat scans for a string checking whether its a float
func ScanStringCheckFloat() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	_, err := utils.ToFloat(inputString)
	if err != nil {
		return "", errors.New("Amount entered is not a float, quitting")
	}
	return inputString, nil
}

// ScanPassword scans for a password
func ScanPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Println(err)
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
		log.Println(err)
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}
