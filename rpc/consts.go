package rpc

import (
	"time"
)

var TimeoutVal = 5 * time.Second

func SetConsts(timeout int) {
	TimeoutVal = time.Duration(time.Duration(timeout) * time.Second)
}
