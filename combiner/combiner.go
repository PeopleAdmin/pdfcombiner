// Package combiner contains methods to consume Job objects, orchestrating
// the download, combination, and upload of a group of related PDF files.
package combiner

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/PeopleAdmin/pdfcombiner/notifier"
	"log"
	"time"
)

var downloadTimeout = 3 * time.Minute

// Combine is the entry point to this package.  Given a Job, downloads
// all the files, combines them into a single one, uploads it to AWS
// and posts the status to a callback endpoint.
func Combine(j *job.Job) {
	defer cleanup(j)
	throttle()
	saveDir := mkTmpDir()
	getAllFiles(j, saveDir)
	if !j.HasDownloadedDocs() {
		return
	}

	componentPaths := fsPathsOf(j.Downloaded, saveDir)
	combinedPath := saveDir + "combined.pdf"
	merger := cpdf.New(combinedPath)
	err := merger.Merge(componentPaths, j.Title)
	if err != nil {
		j.AddError(fmt.Errorf("while combining file: %v", err))
		return
	}
	j.UploadCombinedFile(combinedPath)
	return
}

// Send an update on the success or failure of the operation to the
// callback URL provided by the job originator.
func cleanup(j *job.Job) {
	log.Println("work complete, posting status to callback:", j.Callback)
	notifier.SendNotification(j)
	removeJob()
}
