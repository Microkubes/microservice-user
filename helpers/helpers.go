package helpers

import "time"

func CurrentTimeMilliseconds() int64 {
	return time.Now().UnixNano() / 1000000
}
