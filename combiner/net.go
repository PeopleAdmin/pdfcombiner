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
func getFile(doc *job.Document, dir string, c chan<- *job.Document, e chan<- error) {
	data, err := doc.Get()
	if err != nil {
		e <- fmt.Errorf("%v: %v", doc.Key, err)
		return
	}
	path := localPath(dir, doc)
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		e <- fmt.Errorf("%v: %v", doc.Key, err)
		return
	}
	doc.GetMetadata(cpdf.New(path))
	c <- doc
}

// Fan out workers to download each document in parallel, then block
// until all downloads are complete.
func getAllFiles(j *job.Job, dir string) {
	complete := make(chan *job.Document, j.DocCount())
	failures := make(chan error, j.DocCount())
	for _, doc := range j.DocList {
		go getFile(&doc, dir, complete, failures)
	}

	waitForDownloads(j, complete, failures)
}

// Listen on several channels for information from background download
// tasks -- each task will either send a s.Stat through c, an error through
// e, or timeout.  Once all docs are accounted for, return the total number
// of bytes recieved.
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
