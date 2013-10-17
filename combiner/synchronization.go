package combiner

import (
	"sync/atomic"
	"time"
)

var (
	maxJobs        int = 15
	jobCounter     int32 = 0
	waitingCounter int32 = 0
	jobsReceivedCount int32 = 0
	jobsFinishedCount int32 = 0
)

func CurrentJobs() int {
	return int(atomic.LoadInt32(&jobCounter))
}

func CurrentWait() int {
	return int(atomic.LoadInt32(&waitingCounter))
}

func TotalJobsReceived() int {
	return int(atomic.LoadInt32(&jobsReceivedCount))
}

func TotalJobsFinished() int {
	return int(atomic.LoadInt32(&jobsFinishedCount))
}

// restricts the maximum number of in-progress jobs to prevent the system from
// being overwhelmed with work.
func waitInQueue() {
	addWaiter()
	for CurrentJobs() >= maxJobs {
		time.Sleep(100 * time.Millisecond)
	}
	startJob()
}

func addWaiter() {
	atomic.AddInt32(&waitingCounter, 1)
	atomic.AddInt32(&jobsReceivedCount, 1)
}

func startJob() {
	atomic.AddInt32(&waitingCounter, -1)
	atomic.AddInt32(&jobCounter, 1)
}

func removeJob() {
	atomic.AddInt32(&jobCounter, -1)
	atomic.AddInt32(&jobsFinishedCount, 1)
}

func resetCounters() {
	atomic.StoreInt32(&jobCounter, 0)
	atomic.StoreInt32(&waitingCounter, 0)
}
