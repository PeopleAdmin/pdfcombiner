package job

import (
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"io/ioutil"
	"launchpad.net/goamz/aws"
	"launchpad.net/goamz/s3"
)

// Connect to AWS and add the handle to the Job object.
func (j *Job) connect() (err error) {
	auth, err := aws.EnvAuth()
	if err != nil {
		return
	}
	s := s3.New(auth, aws.USEast)
	j.bucket = s.Bucket(j.BucketName)
	return
}

// Construct an absolute path within a bucket to a given document.
func (j *Job) s3Path(docname string) string {
	return docname
}

// Get retrieves the requested document from S3 as a byte slice or
// decodes it from the embedded `data` attribute of the Document.
func (j *Job) Get(doc Document) (docContent []byte, err error) {
	if testmode.IsEnabled() {
		return
	}
	if doc.Data == "" {
		docContent, err = j.bucket.Get(j.s3Path(doc.Key))
	} else {
		docContent, err = decodeEmbeddedData(doc.Data)
	}
	return
}

// UploadCombinedFile sends a file to the job's CombinedKey on S3.
func (j *Job) UploadCombinedFile(localPath string) (err error) {
	if testmode.IsEnabled() {
		return
	}
	content, err := ioutil.ReadFile(localPath)
	if err != nil {
		j.AddError(j.CombinedKey, err)
		return
	}
	err = j.bucket.Put(j.CombinedKey, content, "application/pdf", s3.Private)
	if err != nil {
		j.AddError(j.CombinedKey, err)
	}
	j.uploadComplete = true
	return
}
