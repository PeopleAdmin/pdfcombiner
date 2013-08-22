package job

import (
	"github.com/brasic/pdfcombiner/stat"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var validJSON = `{"bucket_name":"asd","employer_id":606,"doc_list":[{"name":"100001.pdf"}], "callback":"http://localhost:9090"}`

func TestValidation(t *testing.T) {
	invalidJob := &Job{}
	assert.False(t, invalidJob.IsValid(), "Job should not validate when inputs are invalid")
	newJob := &Job{
		Callback:   "http://blah.com",
		BucketName: "A",
		EmployerID: 1,
		DocList:    make([]Document, 0),
	}
	assert.False(t, newJob.IsValid(), "Job should not validate with an empty doclist")
	newJob.DocList = append(newJob.DocList, Document{Name: "doc.pdf"})
	assert.True(t, newJob.IsValid(), "Job should validate when inputs are valid")
}

func TestJSONDeserialization(t *testing.T) {
	var cases = []struct {
		in       string
		validity bool
	}{
		{"{", false},
		{"{}", false},
		{validJSON, true},
	}
	for _, c := range cases {
		newJob, err := newFromString(c.in)
		assert.Equal(t, newJob.IsValid(), c.validity, "validity should be as expected for "+c.in)
		assert.Equal(t, (err == nil), c.validity, "error should be returned if something went wrong for "+c.in)
	}
}

func TestCompletion(t *testing.T) {
	j, _ := newFromString(validJSON)
	assert.Equal(t, j.CompleteCount(), 0, "New jobs have no complete docs")
	assert.False(t, j.HasDownloadedDocs(), "HasDownloadedDocs() is false when appropriate")
	j.MarkComplete("100001.pdf", stat.Stat{})
	assert.Equal(t, j.CompleteCount(), 1, "Complete count is accurate")
	assert.True(t, j.HasDownloadedDocs(), "HasDownloadedDocs() is true when appropriate")
}

func newFromString(in string) (newJob *Job, err error) {
	json := strings.NewReader(in)
	newJob, err = NewFromJSON(json)
	return
}
