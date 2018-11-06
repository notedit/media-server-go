package sdp

import (
	"fmt"
	"strings"
)

func uintArrayToString(a []uint, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}

func intArrayToString(a []int, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}
