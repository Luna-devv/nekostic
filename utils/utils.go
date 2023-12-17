package utils

import (
	"fmt"
	"strconv"
	"time"
)

func Make6DigitDay() string {
	now := time.Now()

	return fmt.Sprintf("%02d", int64(now.Day())) + strconv.FormatInt(int64(now.Month()), 10) + strconv.FormatInt(int64(now.Year()), 10)
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / 1e6
}
