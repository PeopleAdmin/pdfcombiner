package combiner

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"math/rand"
	"os"
	"strings"
	"time"
)

var basedir = "/tmp/"

// Make and return a randomized temporary directory.
func mkTmpDir() (dirname string) {
	rand.Seed(time.Now().UnixNano())
	dirname = fmt.Sprintf("/tmp/pdfcombiner/%d/", rand.Int())
	os.MkdirAll(dirname, 0777)
	return
}

// Get the absolute paths to a list of docs.
func fsPathsOf(docs []*job.Document, dir string) (paths []string) {
	paths = make([]string, len(docs))
	for idx, doc := range docs {
		paths[idx] = localPath(dir, doc)
	}
	return
}

// localPath replaces any s3 key directory markers with underscores so
// we don't need to recursively create directories when saving files.
func localPath(dir string, doc *job.Document) string {
	return dir + strings.Replace(doc.Key, "/", "_", -1)
}
