package job

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var tmpDirPattern = "/tmp/pdfcombiner/%d/"

func (j *Job) LocalPath() string {
	return j.workingDirectory + "combined.pdf"
}

// Make a randomized temporary directory to work in.
// Panics if it can't be created, since we can't continue.
func (j *Job) mkTmpDir() {
	rand.Seed(time.Now().UnixNano())
	j.workingDirectory = fmt.Sprintf(tmpDirPattern, rand.Int())
	err := os.MkdirAll(j.workingDirectory, 0777);
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
	paths = make([]string, len(j.Downloaded))
	for idx, doc := range j.Downloaded {
		paths[idx] = doc.LocalPath()
	}
	return
}
