package testmode

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var forceState int
var notForced = 0
var forcedTrue = 1
var forcedFalse = 2

var forceToggle os.Signal = syscall.SIGUSR1

func SetupSignals() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, forceToggle)
	go func() {
		for _ = range sigs {
			toggleForce()
		}
	}()
}

func toggleForce() {
	forceState = (forceState + 1) % 3
}

func IsEnabled() bool {
	if !IsEnabledFast() {
		return false
	}
	RandomDelay(3)
	return true
}

func IsEnabledFast() bool {
	return os.Getenv("TEST_MODE") != ""
}

// Returns false with probability P(falsePct)
func RandomBoolean(falsePct float64) bool {
	if forceState == forcedTrue {
		return true
	}
	if forceState == forcedFalse {
		return false
	}
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
