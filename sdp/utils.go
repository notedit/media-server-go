package sdp

import (
	"fmt"
	"strings"
)

func arrayToString(a []int, delim string) string {
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
}
