package combiner

import (
	"errors"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	. "github.com/PeopleAdmin/pdfcombiner/stat"
	"io/ioutil"
	"log"
	"time"
)

// Get an individual file from S3.  If successful, writes the file out to disk
// and sends a stat object back to the main channel.  If there are errors they
// are sent back through the error channel.
func getFile(j *job.Job, doc job.Document, dir string, c chan<- Stat, e chan<- Stat) {
	start := time.Now()
	data, err := j.Get(doc)
	if err != nil {
		e <- Stat{Filename: doc.Key, Err: err}
		return
	}
	path := localPath(dir, doc.Key)
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		e <- Stat{Filename: doc.Key, Err: err}
		return
	}
	pagecount := cpdf.PageCount(path)
	c <- Stat{Filename: doc.Key,
		Size:      len(data),
		PageCount: pagecount,
		DlTime:    time.Since(start)}
}

// Fan out workers to download each document in parallel, then block
// until all downloads are complete.
func getAllFiles(j *job.Job, dir string) {
	start := time.Now()
	complete := make(chan Stat, j.DocCount())
	failures := make(chan Stat, j.DocCount())
	for _, doc := range j.DocList {
		go getFile(j, doc, dir, complete, failures)
	}

	totalBytes := waitForDownloads(j, complete, failures)
	printSummary(start, totalBytes, j.CompleteCount())
}

// Listen on several channels for information from background download
// tasks -- each task will either send a s.Stat through c, an error through
// e, or timeout.  Once all docs are accounted for, return the total number
// of bytes recieved.
func waitForDownloads(j *job.Job, complete, failures <-chan Stat) (totalBytes int) {
	for _ = range j.DocList {
		select {
		case packet := <-complete:
			log.Printf("%s was %d bytes\n", packet.Filename, packet.Size)
			totalBytes += packet.Size
			j.MarkComplete(packet.Filename, packet)
		case bad := <-failures:
			j.AddError(bad.Filename, bad.Err)
		case <-time.After(downloadTimeout):
			j.AddError("general", errors.New("timed out while downloading"))
			return
		}
	}
	return
}
