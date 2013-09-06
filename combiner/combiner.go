// Package combiner contains methods to consume Job objects, orchestrating
// the download, combination, and upload of a group of related PDF files.
package combiner

import (
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/PeopleAdmin/pdfcombiner/notifier"
	. "github.com/PeopleAdmin/pdfcombiner/stat"
	"log"
	"time"
)

var downloadTimeout = 3 * time.Minute

// Combine is the entry point to this package.  Given a Job, downloads
// all the files, combines them into a single one, uploads it to AWS
// and posts the status to a callback endpoint.
func Combine(j *job.Job) {
	defer cleanup(j)
	startAll := time.Now()
	throttle()
	saveDir := mkTmpDir()
	getAllFiles(j, saveDir)
	if !j.HasDownloadedDocs() {
		return
	}

	componentPaths := fsPathsOf(j.Downloaded, saveDir)
	combinedPath := saveDir + "combined.pdf"
	err := cpdf.Merge(componentPaths, combinedPath, j.Title)
	if err != nil {
		j.AddError("general", err)
		return
	}
	startUpload := time.Now()
	j.UploadCombinedFile(combinedPath)
	j.PerfStats["upload"] = Stat{DlTime: time.Since(startUpload)}
	j.PerfStats["total"] = Stat{DlTime: time.Since(startAll)}
	return
}

// Print a summary of the download activity.
func printSummary(start time.Time, bytes int, count int) {
	seconds := time.Since(start).Seconds()
	mbps := float64(bytes) / 1024 / 1024 / seconds
	log.Printf("got %d bytes over %d files in %f secs (%f MB/s)\n",
		bytes, count, seconds, mbps)
}

// Send an update on the success or failure of the operation to the
// callback URL provided by the job originator.
func cleanup(j *job.Job) {
	log.Println("work complete, posting status to callback:", j.Callback)
	notifier.SendNotification(j)
	removeJob()
}
