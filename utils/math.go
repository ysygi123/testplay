package utils

import (
	"math/rand"
	"time"
)

func RangeRand(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
