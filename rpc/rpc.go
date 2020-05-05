package rpc

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// SetupBasicHandlers sets up two handler functions that serve ping and default response at /
func SetupBasicHandlers() {
	SetupDefaultHandler()
	SetupPingHandler()
}

// CheckOrigin checks if the origin of the incoming request is localhost
func CheckOrigin(w http.ResponseWriter, r *http.Request) error {
	if !strings.Contains(r.Header.Get("Origin"), "localhost") {
		ResponseHandler(w, StatusNotFound)
		log.Println("origin not localhost")
		return errors.New("origin not localhost")
	}
	return nil
}

// CheckGet checks if the incoming request is a GET request
func CheckGet(w http.ResponseWriter, r *http.Request) error {
	//err := CheckOrigin(w, r)
	if r.Method != "GET" {
		ResponseHandler(w, StatusNotFound)
		log.Println("method not get")
		return errors.New("method not get")
	}
	return nil
}

// CheckPost checks whether the incoming request is a POST request
func CheckPost(w http.ResponseWriter, r *http.Request) error {
	//err := CheckOrigin(w, r)
	if r.Method != "POST" {
		ResponseHandler(w, StatusNotFound)
		log.Println("method not post")
		return errors.New("method not post")
	}
	return nil
}

// CheckPut checks whether the incoming request is a PUT request
func CheckPut(w http.ResponseWriter, r *http.Request) error {
	err := CheckOrigin(w, r)
	if err != nil || r.Method != "PUT" {
		ResponseHandler(w, StatusNotFound)
		log.Println("method not put or origin not localhost")
		return errors.New("method not put or origin not localhost")
	}
	return nil
}

// GetRequest is a handler that makes it easy to send out GET requests
func GetRequest(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: TimeoutVal,
	}

	res, err := client.Get(url)
	if err != nil {
		res, err = client.Get(url)
		if err != nil {
			log.Println("did not make request: ", err)
			return nil, errors.Wrap(err, "did not make request")
		}
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(res.Body)
}

// PutRequest is a handler that makes it easy to send out PUT requests
func PutRequest(body string, payload io.Reader) ([]byte, error) {
	// the body must be the param that you usually pass to curl's -d option
	client := &http.Client{
		Timeout: TimeoutVal,
	}

	req, err := http.NewRequest("PUT", body, payload)
	if err != nil {
		log.Println("did not create new PUT request: ", err)
		return nil, errors.Wrap(err, "did not create new PUT request")
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request: ", err)
		return nil, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(res.Body)
}

// PostRequest is a handler that makes it easy to send out POST requests
func PostRequest(body string, payload io.Reader) ([]byte, error) {
	// the body must be the param that you usually pass to curl's -d option
	client := &http.Client{
		Timeout: TimeoutVal,
	}

	req, err := http.NewRequest("POST", body, payload)
	if err != nil {
		log.Println("did not create new POST request: ", err)
		return nil, errors.Wrap(err, "did not create new POST request")
	}

	res, err := client.Do(req)
	if err != nil {
		log.Println("did not make request: ", err)
		return nil, errors.Wrap(err, "did not make request")
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	return ioutil.ReadAll(res.Body)
}

// PostForm is a handler that makes it easy to send out POST form requests
func PostForm(body string, postdata url.Values) ([]byte, error) {

	data, err := http.PostForm(body, postdata)
	if err != nil {
		log.Println("could not relay get request: ", err)
		return nil, errors.Wrap(err, "could not relay get request")
	}

	defer func() {
		if ferr := data.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	rData, err := ioutil.ReadAll(data.Body)
	if err != nil {
		log.Println(err)
	}
	return rData, err
}

// GetAndSendJson is a handler that makes a get request and returns json data
func GetAndSendJson(w http.ResponseWriter, body string, x interface{}) {
	data, err := GetRequest(body)
	if err != nil {
		log.Println("did not get response: ", err)
		ResponseHandler(w, StatusBadRequest)
		return
	}
	// now data is in byte, we need the other structure now
	err = json.Unmarshal(data, &x)
	if err != nil {
		log.Println("did not unmarshal json: ", err)
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
		log.Println("did not get response: ", err)
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
		log.Println(err)
		return false, errors.Wrap(err, "could not parse url string")
	}
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		log.Println(err)
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
			} else {
				log.Println(err)
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
		if err := http.ListenAndServe(":80", http.HandlerFunc(HTTPSRedirect)); err != nil {
				log.Fatalf("ListenAndServe error: %v", err)
		}
	*/
}
