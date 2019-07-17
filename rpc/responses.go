package rpc

import (
	"encoding/json"
	//"log"
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

// ResponseHandler is the default response handler that sends out response codes on successful
// completion of certain calls
func ResponseHandler(w http.ResponseWriter, status int) {
	var response StatusResponse
	response.Code = status
	switch status {
	case StatusOK:
		response.Status = "OK"
	case StatusCreated:
		response.Status = "Method Created"
	case StatusMovedPermanently:
		response.Status = "Endpoint moved permanently"
	case StatusBadRequest:
		response.Status = "Bad Request error!"
	case StatusUnauthorized:
		response.Status = "You are unauthorized to make this request"
	case StatusPaymentRequired:
		response.Status = "Payment required before you can access this endpoint"
	case StatusNotFound:
		response.Status = "404 Error Not Found!"
	case StatusInternalServerError:
		response.Status = "Internal Server Error"
	case StatusLocked:
		response.Status = "Endpoint locked until further notice"
	case StatusTooManyRequests:
		response.Status = "Too many requests made, try again later"
	case StatusBadGateway:
		response.Status = "Bad Gateway Error"
	case StatusServiceUnavailable:
		response.Status = "Service Unavailable error"
	case StatusGatewayTimeout:
		response.Status = "Gateway Timeout Error"
	case StatusNotAcceptable:
		response.Status = "Not accepted"
	default:
		response.Status = "404 Page Not Found"
	}
	MarshalSend(w, response)
}

// MarshalSend marshals and writes a json string
func MarshalSend(w http.ResponseWriter, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		WriteToHandler(w, []byte("did not marshal json"))
		return
	}
	WriteToHandler(w, xJson)
}

// WriteToHandler returns a reply to the passed writer
func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
