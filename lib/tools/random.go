package tools

import (
	"math/rand"
)

func GetRandomInt(start, end int) int {
	return rand.Intn(end-start) + start
}

func GetRandomInt64(start, end int64) int64 {
	return rand.Int63n(end-start) + start
}
