package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// GetRequest is a handler that makes it easy to send out GET requests
// we don't set timeouts here because block times can be variable and a single request
// can sometimes take a long while to complete
func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("did not create new GET request", err)
		return dummy, err
	}
	req.Header.Set("Origin", "localhost")
	req.SetBasicAuth(APIKey, APISecret)
	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// PutRequest is a handler that makes it easy to send out PUT requests
func PutRequest(body string, payload io.Reader) ([]byte, error) {

	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("PUT", body, payload)
	if err != nil {
		log.Println("did not create new PUT request", err)
		return dummy, err
	}
	// need to add this header or we'll get a negative response
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("did not read from ioutil", err)
		return dummy, err
	}

	return x, nil
}

// PostRequest is a handler that makes it easy to send out POST requests
func PostRequest(body string, payload io.Reader) ([]byte, error) {

	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("POST", body, payload)
	if err != nil {
		log.Println("did not create new POST request", err)
		return dummy, err
	}

	req.SetBasicAuth(APIKey, APISecret)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("did not make request", err)
		return dummy, err
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("did not read from ioutil", err)
		return dummy, err
	}

	return x, nil
}
