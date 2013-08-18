package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailingValidation(t *testing.T) {
	newJob := &Job{}
	assert.False(t, newJob.IsValid())
}

func TestPassingValidation(t *testing.T) {
	newJob := &Job{
		Callback: "A", BucketName: "A", EmployerId: 1, DocList: []string{"A"}}
	assert.True(t, newJob.IsValid())
}
