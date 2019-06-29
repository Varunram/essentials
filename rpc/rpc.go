package rpc

// the rpc package contains functions related to the server which will be interacting
// with the frontend. Not expanding on this too much since this will be changing quite often
// also evaluate on how easy it would be to rewrite this in nodeJS since the
// frontend is in react. Not many advantages per se and this works fine, so I guess
// we'll stay with this one for a while
import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// API documentation over at the apidocs repo

// StatusResponse defines a generic status response structure
type StatusResponse struct {
	Code   int
	Status string
}

// setupBasicHandlers sets up two handler functions that can be used to serve a default
// 404 response when we either error out or received input is incorrect.  This is not
// exactly ideal, because we don't expcet the RPC to be exposed and would like some more
// errors when we handle it on the frontend, but this makes for more a bit more
// secure Frontedn implementation which doesn't leak any information to the frontend
func SetupBasicHandlers() {
	SetupDefaultHandler()
	SetupPingHandler()
}

// CheckOrigin checks the origin of the incoming request
func CheckOrigin(w http.ResponseWriter, r *http.Request) {
	// re-enable this function for all private routes
	if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
		http.Error(w, "404 page not found", http.StatusNotFound)
	}
}

// CheckGet checks if the invoming request is a GET request
func CheckGet(w http.ResponseWriter, r *http.Request) {
	CheckOrigin(w, r)
	if r.Method != "GET" {
		ResponseHandler(w, StatusNotFound)
		return
	}
}

// checkPost checks whether the incomign request is a POST request
func CheckPost(w http.ResponseWriter, r *http.Request) {
	CheckOrigin(w, r)
	if r.Method != "POST" {
		ResponseHandler(w, StatusNotFound)
		return
	}
}

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

// setupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func SetupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		CheckGet(w, r)
		CheckOrigin(w, r)
		ResponseHandler(w, StatusNotFound)
	})
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func SetupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		CheckGet(w, r)
		CheckOrigin(w, r)
		ResponseHandler(w, StatusOK)
	})
}
