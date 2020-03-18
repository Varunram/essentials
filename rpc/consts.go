package rpc

import (
	"time"
)

// TimeoutVal is the timeout associated with a single call
var TimeoutVal = 5 * time.Second

// SetConsts is used to set the timeout interval for requests made
func SetConsts(timeout int) {
	TimeoutVal = time.Duration(time.Duration(timeout) * time.Second)
}
