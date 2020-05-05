package rpc

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

// SetupLocalHttpsClient can be used to setup a local client configured to accept a user generated cert
func SetupLocalHttpsClient(path string, timeout time.Duration) *http.Client {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("failed to read from file: ", err)
		panic(err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	config := &tls.Config{
		RootCAs: rootCAs,
	}

	tr := &http.Transport{TLSClientConfig: config}
	return &http.Client{Transport: tr, Timeout: timeout}
}

// HttpsGet is a function that should only be used on localhost with a client configured to accept a user generated cert
func HttpsGet(client *http.Client, url string) ([]byte, error) {
	// Read in the cert file
	res, err := client.Get(url)
	if err != nil {
		log.Println("did not make request: ", err)
		return nil, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	return data, err
}

// HttpsPost is a function that should only be used on localhost with a client configured to accept a user generated cert
func HttpsPost(client *http.Client, url string, postdata url.Values) ([]byte, error) {
	// Read in the cert file
	res, err := client.PostForm(url, postdata)
	if err != nil {
		log.Println("did not make request: ", err)
		return nil, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	return data, err
}
