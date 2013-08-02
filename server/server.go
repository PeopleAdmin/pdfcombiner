package pdfcombiner

import (
  "pdfcombiner/combiner"
  "net/http"
  "net/url"
  "fmt"
)

type CombinerServer struct {}

// Validate that the required URL params are present and correct.
func check_params(params map[string] []string ) (docs []string, callback string, ok bool) {
  callbacks := params["callback"]
  if len(callbacks) < 1 {
    return
  }
  callback = callbacks[0]
  _, parseErr := url.Parse(callback)
  docs, docsPresent := params["docs"]
  ok = docsPresent && (parseErr == nil)
  return
}

// Looks for one or more ?docs=FILE params and if found starts combination.
func (c CombinerServer) ProcessCombineRequest(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  params := r.Form
  docs, callback, ok := check_params(params)
  if !ok {
    http.Error(w, "Need some docs and a callback url", http.StatusBadRequest)
    return
  }
  request := &pdfcombiner.CombineRequest{
    BucketName: "pa-hrsuite-production",
    EmployerId: "606",
    DocList: docs,
    Callback: callback }
  go pdfcombiner.Combine(request)
  fmt.Fprintln(w, "Started combination on",docs)
}

func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {}
