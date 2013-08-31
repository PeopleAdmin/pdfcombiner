package testmode

import (
	"math/rand"
	"os"
	"time"
)

func IsEnabled() bool {
	if !IsEnabledFast() {
		return false
	}
	time.Sleep(1 * time.Second)
	return true
}

func IsEnabledFast() bool {
	return os.Getenv("TEST_MODE") != ""
}

// Returns false with probability P(falsePct)
func RandomBoolean(falsePct float64) bool {
	rand.Seed(time.Now().UnixNano())
	if falsePct > rand.Float64() {
		return false
	}
	return true
}
