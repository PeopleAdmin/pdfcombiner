package server

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	handler        = new(CombinerServer)
	invalidJson    = "{"
	incompleteJson = "{}"
	validJson      = "{\"bucket_name\":\"a\",\"employer_id\":1,\"doc_list\":[\"1.pdf\"], \"callback\":\"http://a.com\"}"
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
	postBody := strings.NewReader(invalidJson)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(invalidMessage)+"\n")
}

func TestRejectsIncompleteRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	postBody := strings.NewReader(incompleteJson)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(invalidMessage)+"\n")
}

func TestAcceptsValidRequest(t *testing.T) {
	recorder := httptest.NewRecorder()
	postBody := strings.NewReader(validJson)
	req, _ := http.NewRequest("POST", "http://example.com/", postBody)
	handler.ProcessJob(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, recorder.Body.String(), string(okMessage))
}
