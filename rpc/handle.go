package rpc

import (
	"log"
	"net/http"
)

// Err prints an error and returns to the caller
func Err(w http.ResponseWriter, err error, status int, msgs ...string) bool {
	var retmsg bool
	var print bool

	if len(msgs) > 0 {
		print = true
	}

	if len(msgs) == 2 {
		retmsg = true
	}

	if err != nil {
		if print {
			log.Println(err, msgs[0])
		}
		if retmsg {
			ResponseHandler(w, StatusBadRequest, msgs[1])
		} else {
			ResponseHandler(w, status)
		}
		return true
	}

	return false
}
