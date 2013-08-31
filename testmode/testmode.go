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
	RandomDelay(8)
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

// Waits for up to maxWait seconds
func RandomDelay(maxWait int64) {
	rand.Seed(time.Now().UnixNano())
	waitRange := maxWait * int64(time.Second)
	waitTime := rand.Int63n(waitRange)
	time.Sleep(time.Duration(waitTime))
}
