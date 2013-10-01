// Package combiner contains methods to consume Job objects, orchestrating
// the download, combination, and upload of a group of related PDF files.
package combiner

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/PeopleAdmin/pdfcombiner/notifier"
	"log"
	"strings"
	"time"
)

var downloadTimeout = 3 * time.Minute

// Combine is the entry point to this package.  Given a Job, downloads
// all the files, combines them into a single one, uploads it to AWS
// and posts the status to a callback endpoint.
func Combine(j *job.Job) {
	defer cleanup(j)
	waitInQueue()
	getAllFiles(j)
	if !j.HasDownloadedDocs() {
		return
	}

	j.SortDownloaded()
	j.DownloadsDoneAt = time.Now()
	j.CleanupBrokenDocs()
	err := cpdf.Merge(j)
	if err != nil {
		j.AddError(fmt.Errorf("while combining file: %v", err))
		return
	}
	j.UploadCombinedFile()
	return
}

// Send an update on the success or failure of the operation to the
// callback URL provided by the job originator.
func cleanup(j *job.Job) {
	logErrors(j)
	notifier.SendNotification(j)
	removeJob()
	j.Cleanup()
}

func logErrors(j *job.Job) {
	if len(j.Errors) > 0 {
		log.Printf("ERRORS (%d) encountered while processing '%s':\n", len(j.Errors), j.Callback)
		for i, err := range j.Errors {
			log.Printf("  %d: %v", i+1, strings.Replace(err.Error(), "\n", `\n`, -1))
		}
	}
}
