package job

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

var (
	awsRegion              = aws.USEast
	uploadedFilePermission = s3.Private
)

type Downloadable interface {
	s3Path() string
}

// Get retrieves the requested document, either from S3 or by decoding the
// embedded `Data` attribute of the Document.
func (doc *Document) Get() (docContent []byte, err error) {
	if testmode.IsEnabled() {
		return
	}
	switch doc.Data {
	case "":
		return doc.parent.download(doc)
	default:
		return decodeEmbeddedData(doc.Data)
	}
}

// UploadCombinedFile sends a file to the job's CombinedKey on S3.
func (j *Job) UploadCombinedFile() (err error) {
	if testmode.IsEnabled() {
		return
	}
	content, err := ioutil.ReadFile(j.LocalPath())
	if err != nil || len(content) == 0 {
		j.AddError(fmt.Errorf("Reading file '%v' for upload: %v", j.LocalPath(), err))
		return
	}
	err = j.bucket.Put(j.CombinedKey, content, "application/pdf", uploadedFilePermission)
	if err != nil {
		j.AddError(fmt.Errorf("Uploading to S3 key '%s': %v", j.CombinedKey, err))
	}
	j.uploadComplete = true
	return
}

// Connect to AWS and add the handle to the Job object.
func (j *Job) connect() (err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	s := s3.New(auth, awsRegion)
	j.bucket = s.Bucket(j.BucketName)
	return
}

func (j *Job) download(doc Downloadable) (content []byte, err error) {
	remotePath := doc.s3Path()
	return j.bucket.Get(remotePath)
}
