package scan

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/pkg/errors"

	utils "github.com/Varunram/essentials/utils"
	"golang.org/x/crypto/ssh/terminal"
)

// package scan can be used by CLI clients that want to accept inptus from the CLI

// Int scans for an integer
func Int() (int, error) {
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

// Float scans for a float
func Float() (float64, error) {
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

// String scans for a string
func String() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", errors.New("Couldn't read user input")
	}
	inputString := scanner.Text()
	return inputString, nil
}

// StringCheckInt scans for a string checking whether it is an integer
func StringCheckInt() (string, error) {
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

// StringCheckFloat scans for a string checking whether its a float
func StringCheckFloat() (string, error) {
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

// Password scans for a password
func Password() (string, error) {
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

// RawPassword scans for a raw password
func RawPassword() (string, error) {
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		log.Println(err)
		return "", err
	}
	password := string(bytePassword)
	return password, nil
}
