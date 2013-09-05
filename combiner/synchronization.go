package combiner

import (
	"sync/atomic"
	"time"
)

var (
	maxJobs        int32 = 15
	jobCounter     int32 = 0
	waitingCounter int32 = 0
)

func CurrentJobs() int32 {
	return atomic.LoadInt32(&jobCounter)
}

func CurrentWait() int32 {
	return atomic.LoadInt32(&waitingCounter)
}

// throttle prevents the system from being overwhelmed with work.
// Blocks until the number of active jobs is less than a preset threshold.
func throttle() {
	addWaiter()
	for CurrentJobs() >= maxJobs {
		time.Sleep(100 * time.Millisecond)
	}
	addJob()
}

func addWaiter() {
	atomic.AddInt32(&waitingCounter, 1)
}

func addJob() {
	atomic.AddInt32(&waitingCounter, -1)
	atomic.AddInt32(&jobCounter, 1)
}

func removeJob() {
	atomic.AddInt32(&jobCounter, -1)
}
