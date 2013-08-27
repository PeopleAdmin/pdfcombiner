// Package combiner contains methods to consume Job objects, orchestrating
// the download, combination, and upload of a group of related PDF files.
// TODO now that s.Stat{} can return errors, maybe the err chan is unnecessary
package combiner

import (
	"errors"
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/cpdf"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/PeopleAdmin/pdfcombiner/notifier"
	s "github.com/PeopleAdmin/pdfcombiner/stat"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	downloadTimeout = 3 * time.Minute
	basedir         = "/tmp/"
	maxGoroutines   = 30
)

// Get an individual file from S3.  If successful, writes the file out to disk
// and sends a stat object back to the main channel.  If there are errors they
// are sent back through the error channel.
func getFile(j *job.Job, docname string, dir string, c chan<- s.Stat, e chan<- s.Stat) {
	start := time.Now()
	data, err := j.Get(docname)
	if err != nil {
		e <- s.Stat{Filename: docname, Err: err}
		return
	}
	path := localPath(dir, docname)
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		e <- s.Stat{Filename: docname, Err: err}
		return
	}
	pagecount := cpdf.PageCount(path)
	c <- s.Stat{Filename: docname,
		Size:      len(data),
		PageCount: pagecount,
		DlTime:    time.Since(start)}
}

// Fan out workers to download each document in parallel, then block
// until all downloads are complete.
func getAllFiles(j *job.Job, dir string) {
	start := time.Now()
	c := make(chan s.Stat, j.DocCount())
	e := make(chan s.Stat, j.DocCount())
	for _, doc := range j.DocList {
		throttle()
		go getFile(j, doc.Name, dir, c, e)
	}

	totalBytes := waitForDownloads(j, c, e)
	printSummary(start, totalBytes, j.CompleteCount())
}

// Prevents the system from being overwhelmed with work.
// Blocks until the number of Goroutines is less than a preset threshold.
func throttle() {
	for {
		if runtime.NumGoroutine() < maxGoroutines {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// Listen on several channels for information from background download
// tasks -- each task will either send a s.Stat through c, an error through
// e, or timeout.  Once all docs are accounted for, return the total number
// of bytes recieved.
func waitForDownloads(j *job.Job, c <-chan s.Stat, e <-chan s.Stat) (totalBytes int) {
	for _ = range j.DocList {
		select {
		case packet := <-c:
			log.Printf("%s was %d bytes\n", packet.Filename, packet.Size)
			totalBytes += packet.Size
			j.MarkComplete(packet.Filename, packet)
		case bad := <-e:
			j.AddError(bad.Filename, bad.Err)
		case <-time.After(downloadTimeout):
			j.AddError("general", errors.New("timed out while downloading"))
			return
		}
	}
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
func postToCallback(j *job.Job) {
	log.Println("work complete, posting status to callback:", j.Callback)
	notifier.SendNotification(j)
}

// Make and return a randomized temporary directory.
func mkTmpDir() (dirname string) {
	rand.Seed(time.Now().UnixNano())
	dirname = fmt.Sprintf("/tmp/pdfcombiner/%d/", rand.Int())
	os.MkdirAll(dirname, 0777)
	return
}

// Combine is the entry point to this package.  Given a Job, downloads
// all the files, combines them into a single one, uploads it to AWS
// and posts the status to a callback endpoint.
func Combine(j *job.Job) bool {
	defer postToCallback(j)
	saveDir := mkTmpDir()
	getAllFiles(j, saveDir)
	if j.HasDownloadedDocs() {
		err := cpdf.Merge(fsPathsOf(j.Downloaded, saveDir), saveDir+"combined.pdf")
		if err != nil {
			j.AddError("general", err)
			return false
		}
	}
	j.UploadCombinedFile(saveDir + "combined.pdf")
	return true
}

// Get the absolute paths to a list of docs.
func fsPathsOf(docs []string, dir string) (paths []string) {
	paths = make([]string, len(docs))
	for idx, doc := range docs {
		paths[idx] = localPath(dir, doc)
	}
	return
}

// localPath replaces any s3 key directory markers with underscores so
// we don't need to recursively create directories when saving files.
func localPath(dir, remotePath string) string {
	return dir + strings.Replace(remotePath, "/", "_", -1)
}
