package job

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var tmpDirPattern = "/tmp/pdfcombiner/%s.%s/"

func (j *Job) LocalPath() string {
	return j.workingDirectory + "combined.pdf"
}

// Make a temporary directory to work in based on the job dir plus the current
// time.  Panics if it can't be created, since we can't continue.
func (j *Job) mkTmpDir() {
	timeSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[8:14]
	j.workingDirectory = fmt.Sprintf(tmpDirPattern, j.Id(), timeSuffix)
	err := os.MkdirAll(j.workingDirectory, 0777)
	if err != nil {
		panic(err)
	}
	return
}

// Delete the tmp directory for this job.
func (j *Job) Cleanup() {
	if j.workingDirectory != "" && os.Getenv("DEBUG") == "" {
		os.RemoveAll(j.workingDirectory)
	}
}

// Get the absolute paths to the completed docs.
func (j *Job) ComponentPaths() (paths []string) {
	paths = make([]string, len(j.DocList))
	for idx, doc := range j.DocList {
		paths[idx] = doc.LocalPath()
	}
	return
}

// For every doc that failed for whatever reason, replace it with a page
// containing an error message.
func (j *Job) CleanupBrokenDocs() {
	for i, doc := range j.DocList {
		if !j.wasDownloaded(doc) {
			j.DocList[i].writeErrorDocument()
		}
	}
}

// Determine if the given Doc is present in the list of downloaded documents.
// In other words, find out if something went wrong.
func (j Job) wasDownloaded(inputDoc Document) bool {
	for _, doc := range j.Downloaded {
		if inputDoc.Key == doc.Key {
			return true
		}
	}
	return false
}
