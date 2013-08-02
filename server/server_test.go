package pdfcombiner

import (
"net/http"
"net/http/httptest"
"github.com/stretchr/testify/assert"
"testing"
)

var(
  handler = new(CombinerServer)
  validParams = "?docs=1.pdf&callback=http://google.com?_test=true"
  invalidMsg = "Need some docs and a callback url\n"
  validMsg = "Started combination on [1.pdf]\n"
)

func TestPing(t *testing.T) {
  recorder := httptest.NewRecorder()
  req, _ := http.NewRequest("GET", "http://example.com/health_check.html", nil)
	handler.Ping(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "", recorder.Body.String())
}

func TestProcessCombineRequestRejectsInvalidRequest(t *testing.T) {
  recorder := httptest.NewRecorder()
  req, _:= http.NewRequest("GET", "http://example.com/", nil)
	handler.ProcessCombineRequest(recorder, req)
	assert.Equal(t, 400, recorder.Code)
	assert.Equal(t, invalidMsg, recorder.Body.String())
}

func TestProcessCombineRequestAcceptsValidRequest(t *testing.T) {
  recorder := httptest.NewRecorder()
  req, _:= http.NewRequest("GET", "http://example.com/"+validParams, nil)
	handler.ProcessCombineRequest(recorder, req)
	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, validMsg, recorder.Body.String())
}
