package combiner

import (
	"errors"
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"io/ioutil"
	"time"
)

// Get an individual file from S3.  If successful, writes the file out to disk
// and sends a stat object back to the main channel.  If there are errors they
// are sent back through the error channel.
func getFile(doc *job.Document, c chan<- *job.Document, e chan<- error) {
	data, err := doc.Get()
	if err != nil {
		e <- fmt.Errorf("%v: %v", doc.Key, err)
		return
	}
	path := doc.LocalPath()
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		e <- fmt.Errorf("%v: %v", doc.Key, err)
		return
	}
	err = doc.GetMetadata(cpdf.New(path, doc.Id()))
	if err != nil {
		e <- fmt.Errorf("While extracting metadata from '%v': %v", doc.Key, err)
		return
	}
	c <- doc
}

// Fan out workers to download each document in parallel, then block
// until all downloads are complete.
func getAllFiles(j *job.Job) {
	complete := make(chan *job.Document, j.DocCount())
	failures := make(chan error, j.DocCount())
	for i := range j.DocList {
		go getFile(&j.DocList[i], complete, failures)
	}

	waitForDownloads(j, complete, failures)
}

// Listen on several channels for information from background download
// tasks -- each task will either send a s.Stat through c, an error through
// e, or timeout.  Once all docs are accounted for, return the total number
// of bytes received.
func waitForDownloads(j *job.Job, complete <-chan *job.Document, failures <-chan error) {
	for _ = range j.DocList {
		select {
		case done := <-complete:
			j.MarkComplete(done)
		case failed := <-failures:
			j.AddError(failed)
		case <-time.After(downloadTimeout):
			j.AddError(errors.New("timed out while downloading"))
			return
		}
	}
	return
}
