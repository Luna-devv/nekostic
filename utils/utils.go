package utils

import (
	"fmt"
	"time"
)

func Make6DigitDay() string {
	now := time.Now()

	return fmt.Sprintf("%02d", now.Day()) + fmt.Sprintf("%02d", now.Month()) + fmt.Sprintf("%04d", now.Year())
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
