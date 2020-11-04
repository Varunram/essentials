package rpc

import (
	"encoding/json"
	"log"
	"net/http"
)

// define these here since we only have stuff that's needed / supported
const (
	StatusOK                  = http.StatusOK                  //  200 RFC 7231, 6.3.1
	StatusCreated             = http.StatusCreated             //  201 RFC 7231, 6.3.2
	StatusMovedPermanently    = http.StatusMovedPermanently    //  301 RFC 7231, 6.4.2
	StatusBadRequest          = http.StatusBadRequest          //  400 RFC 7231, 6.5.1
	StatusUnauthorized        = http.StatusUnauthorized        //  401 RFC 7235, 3.1
	StatusPaymentRequired     = http.StatusPaymentRequired     //  402 RFC 7231, 6.5.2
	StatusNotFound            = http.StatusNotFound            //  404 RFC 7231, 6.5.4
	StatusInternalServerError = http.StatusInternalServerError //  500 RFC 7231, 6.6.1
	StatusBadGateway          = http.StatusBadGateway          //  502 RFC 7231, 6.6.3
	StatusLocked              = http.StatusLocked              //  423 RFC 4918, 11.3
	StatusTooManyRequests     = http.StatusTooManyRequests     //  429 RFC 6585, 4
	StatusGatewayTimeout      = http.StatusGatewayTimeout      //  504 RFC 7231, 6.6.5
	StatusNotAcceptable       = http.StatusNotAcceptable       //  406 RFC 7231, 6.5.6
	StatusServiceUnavailable  = http.StatusServiceUnavailable  //  503 RFC 7231, 6.6.4
)

// StatusResponse defines a generic status response structure
type StatusResponse struct {
	Code    int
	Status  string
	Message string
}

// ResponseHandler is the default response handler that sends out response codes on successful
// completion of certain calls
func ResponseHandler(w http.ResponseWriter, status int, messages ...string) {
	var response StatusResponse
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	response.Code = status
	switch status {
	case StatusOK:
		w.WriteHeader(StatusOK)
		response.Status = "OK"
	case StatusCreated:
		w.WriteHeader(StatusCreated)
		response.Status = "Method Created"
	case StatusMovedPermanently:
		w.WriteHeader(StatusMovedPermanently)
		response.Status = "Endpoint moved permanently"
	case StatusBadRequest:
		w.WriteHeader(StatusBadRequest)
		response.Status = "Bad Request error!"
	case StatusUnauthorized:
		w.WriteHeader(StatusUnauthorized)
		response.Status = "You are unauthorized to make this request"
	case StatusPaymentRequired:
		w.WriteHeader(StatusPaymentRequired)
		response.Status = "Payment required before you can access this endpoint"
	case StatusNotFound:
		w.WriteHeader(StatusNotFound)
		response.Status = "404 Error Not Found!"
	case StatusInternalServerError:
		w.WriteHeader(StatusInternalServerError)
		response.Status = "Internal Server Error"
	case StatusLocked:
		w.WriteHeader(StatusLocked)
		response.Status = "Endpoint locked until further notice"
	case StatusTooManyRequests:
		w.WriteHeader(StatusTooManyRequests)
		response.Status = "Too many requests made, try again later"
	case StatusBadGateway:
		w.WriteHeader(StatusBadGateway)
		response.Status = "Bad Gateway Error"
	case StatusServiceUnavailable:
		w.WriteHeader(StatusServiceUnavailable)
		response.Status = "Service Unavailable error"
	case StatusGatewayTimeout:
		w.WriteHeader(StatusGatewayTimeout)
		response.Status = "Gateway Timeout Error"
	case StatusNotAcceptable:
		w.WriteHeader(StatusNotAcceptable)
		response.Status = "Not accepted"
	default:
		response.Status = "404 Page Not Found"
	}

	var message string
	if len(messages) > 0 {
		message = messages[0]
	}

	response.Message = message
	MarshalSend(w, response)
}

// MarshalSend marshals and writes a json string
func MarshalSend(w http.ResponseWriter, x interface{}) {
	w.Header().Add("Access-Control-Allow-Headers", "Access-Control-Allow-Origin")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	xJSON, err := json.Marshal(x)
	if err != nil {
		log.Println("could not marshal json: ", err)
		WriteToHandler(w, []byte("did not marshal json"))
		return
	}
	WriteToHandler(w, xJSON)
}

// WriteToHandler returns a reply to the passed writer
func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Write(jsonString)
}
