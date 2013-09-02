package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"sync"
)

var (
	handler        = new(CombinerServer)
	invalidJSON    = "{"
	incompleteJSON = "{}"
	validJSON      = `{"bucket_name":"asd","doc_list":[{"name":"100001.pdf"}], "callback":"http://localhost:9090","combined_key":"out.pdf"}`
)

func TestPing(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://example.com/health_check.html", nil)
	handler.Ping(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "", recorder.Body.String())
}

func TestRejectsInvalidRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	postBody := strings.NewReader(invalidJSON)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(invalidMessage))
}

func TestRejectsIncompleteRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	postBody := strings.NewReader(incompleteJSON)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(invalidMessage))
}

func TestAcceptsValidRequest(t *testing.T) {
	handler.pending = new(sync.WaitGroup)
	recorder := httptest.NewRecorder()
	postBody := strings.NewReader(validJSON)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(okMessage))
}
