package rpc

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// package rpc contains stuff that one would most likely define for their own database

// StatusResponse defines a generic status response structure
type StatusResponse struct {
	Code   int
	Status string
}

// SetupBasicHandlers sets up two handler functions that serve ping and default response at /
func SetupBasicHandlers() {
	SetupDefaultHandler()
	SetupPingHandler()
}

// CheckOrigin checks if the origin of the incoming request is localhost
func CheckOrigin(w http.ResponseWriter, r *http.Request) error {
	if !strings.Contains(r.Header.Get("Origin"), "localhost") {
		return errors.New("origin not localhost")
	}
	return nil
}

// CheckGet checks if the incoming request is a GET request
func CheckGet(w http.ResponseWriter, r *http.Request) error {
	//err := CheckOrigin(w, r)
	if r.Method != "GET" {
		ResponseHandler(w, StatusNotFound)
		return errors.New("method not get or origin not localhost")
	}
	return nil
}

// CheckPost checks whether the incoming request is a POST request
func CheckPost(w http.ResponseWriter, r *http.Request) error {
	//err := CheckOrigin(w, r)
	if r.Method != "POST" {
		ResponseHandler(w, StatusNotFound)
		return errors.New("method not post")
	}
	return nil
}

// CheckPut checks whether the incoming request is a PUT request
func CheckPut(w http.ResponseWriter, r *http.Request) error {
	err := CheckOrigin(w, r)
	if err != nil || r.Method != "PUT" {
		ResponseHandler(w, StatusNotFound)
		return errors.New("method not put or origin not localhost")
	}
	return nil
}

// GetRequest is a handler that makes it easy to send out GET requests
func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{
		Timeout: TimeoutVal,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dummy, errors.Wrap(err, "did not create new GET request")
	}
	req.Header.Set("Origin", "localhost")
	res, err := client.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close() ; ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(res.Body)
}

// PutRequest is a handler that makes it easy to send out PUT requests
func PutRequest(body string, payload io.Reader) ([]byte, error) {
	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("PUT", body, payload)
	if err != nil {
		return dummy, errors.Wrap(err, "did not create new PUT request")
	}
	// need to add this header or we'll get a negative response sometimes
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close() ; ferr != nil {
			err = ferr
		}
	}()

	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dummy, errors.Wrap(err, "did not read from ioutil")
	}

	return x, nil
}

// PostRequest is a handler that makes it easy to send out POST requests
func PostRequest(body string, payload io.Reader) ([]byte, error) {
	// the body must be the param that you usually pass to curl's -d option
	var dummy []byte
	req, err := http.NewRequest("POST", body, payload)
	if err != nil {
		return dummy, errors.Wrap(err, "did not create new POST request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close() ; ferr != nil {
			err = ferr
		}
	}()


	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return dummy, errors.Wrap(err, "did not read from ioutil")
	}

	return x, nil
}

// PostForm is a handler that makes it easy to send out POST form requests
func PostForm(body string, postdata url.Values) ([]byte, error) {

	httpdata, err := http.PostForm(body, postdata)
	if err != nil {
		return nil, errors.Wrap(err, "could not relay get request")
	}

	defer func() {
		if ferr := httpdata.Body.Close() ; ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(httpdata.Body)
}

// GetAndSendJson is a handler that makes a get request and returns json data
func GetAndSendJson(w http.ResponseWriter, body string, x interface{}) {
	data, err := GetRequest(body)
	if err != nil {
		log.Println("did not get response:", err)
		ResponseHandler(w, StatusBadRequest)
		return
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		ResponseHandler(w, StatusInternalServerError)
		return
	}
	MarshalSend(w, x)
}

// GetAndSendByte is a handler that makes a get request and returns byte data. This is used
// in cases for which we don't know the format of the returned data, so we can't parse it
func GetAndSendByte(w http.ResponseWriter, body string) {
	data, err := GetRequest(body)
	if err != nil {
		log.Println("did not get response", err)
		ResponseHandler(w, StatusBadRequest)
		return
	}

	w.Write(data)
}

// PutAndSend is a handler that PUTs data and returns the response
func PutAndSend(w http.ResponseWriter, body string, payload io.Reader) {
	data, err := PutRequest(body, payload)
	if err != nil {
		log.Println("did not receive success response", err)
		ResponseHandler(w, StatusBadRequest)
		return
	}
	var x interface{}
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json", err)
		ResponseHandler(w, StatusInternalServerError)
		return
	}
	MarshalSend(w, x)
}

// SetupDefaultHandler sets up the default handler (ie returns 404 for invalid routes)
func SetupDefaultHandler() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// default to 404 for every application not running on localhost
		err := CheckGet(w, r)
		if err != nil {
			log.Println(err) // don't return a response since we've already written to the API caller
			return
		}
		ResponseHandler(w, StatusNotFound) // default response for routes not found is 404
	})
}

// SetupPingHandler is a ping route for remote callers to check if the platform is up
func SetupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err := CheckGet(w, r)
		if err != nil {
			log.Println(err) // don't return a response since we've already written to the API caller
			return
		}
		ResponseHandler(w, StatusOK)
	})
}

// CheckHTTPSRedirect checks if HTTP requests are redirected to HTTPS on the same host
func CheckHTTPSRedirect(urlString string) (bool, error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return false, errors.Wrap(err, "could not parse url string")
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return false, errors.Wrap(err, "errros while constructing new get request")
	}

	var client http.Client
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	req.Close = true

	resp, err := client.Do(req)
	if err == nil && (resp.StatusCode == 301 || resp.StatusCode == 302) { // 301 moved permanently and 302 moved temporarily
		headers := resp.Header
		response, ok := headers["Location"] // get the location of the redirect
		if ok {
			redirURL, err := url.Parse(response[0]) // parse the redirect URL
			if err == nil {
				// if the host is the same and the URL is https
				if redirURL.Host == url.Host && redirURL.Scheme == "https" {
					return true, nil
				}
			}
		}
	}
	// HTTP does not redirect to HTTPS, or does so improperly
	return false, errors.New("does not redirect to https")
}

// HTTPSRedirect redirects servers to use https instead of http
func HTTPSRedirect(w http.ResponseWriter, r *http.Request, origin string) {
	http.Redirect(w, r, "https://"+origin+r.RequestURI, http.StatusMovedPermanently)
	// use the above like so:
	/*
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
				log.Fatalf("ListenAndServe error: %v", err)
		}
	*/
}
