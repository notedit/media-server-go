package mediaserver

import (
	"fmt"
	"math/rand"
	"time"
)

const ssrcmin = 1000000000

// todo
func GenerateSSRC() uint {
	rand.Seed(time.Now().Unix())
	randNum := rand.Uint32()
	if randNum < ssrcmin {
		randNum = randNum + ssrcmin
	}
	fmt.Printf("rand is %v\n", randNum)
	return uint(randNum)
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
