package testmode

import (
	"os"
	"time"
)

func IsEnabled() bool {
	if os.Getenv("TEST_MODE") == "" {
		return false
	}
	time.Sleep(1 * time.Second)
	return true
}
