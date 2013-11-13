package combiner

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestCounts(t *testing.T) {
	resetCounters()
	assert.Equal(t, CurrentJobs(), 0)
	assert.Equal(t, CurrentWait(), 0)
	addWaiter()
	assert.Equal(t, CurrentJobs(), 0)
	assert.Equal(t, CurrentWait(), 1)
	startJob()
	assert.Equal(t, CurrentJobs(), 1)
	assert.Equal(t, CurrentWait(), 0)
	removeJob()
	assert.Equal(t, CurrentJobs(), 0)
	assert.Equal(t, CurrentWait(), 0)
}

// Ensure the counters can withstand concurrent access without being corrupted.
func TestAtomicAccess(t *testing.T) {
	resetCounters()
	responses := make(chan bool, 1000)
	for i := 0; i < 1000; i++ {
		go func(done chan bool) {
			sleepRand()
			addWaiter()
			sleepRand()
			startJob()
			sleepRand()
			removeJob()
			done <- true
		}(responses)
	}
	for i := 0; i < 1000; i++ {
		<-responses
	}
	assert.Equal(t, CurrentJobs(), 0)
	assert.Equal(t, CurrentWait(), 0)
}

func sleepRand() {
	time.Sleep(time.Duration(rand.Intn(200)) * time.Microsecond)
}
