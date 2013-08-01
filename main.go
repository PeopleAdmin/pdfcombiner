// Benchmark parallel downloading and combination of some pdf files.
package main

import (
  "pdfcombiner/combiner"
  "net/http"
  "net/url"
  "fmt"
)

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
func combineEndpoint(w http.ResponseWriter, r *http.Request) {
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

func noopEndpoint(w http.ResponseWriter, r *http.Request) {}

func main() {
  http.HandleFunc("/favicon.ico", noopEndpoint)
  http.HandleFunc("/health_check.html", noopEndpoint)
  http.HandleFunc("/", combineEndpoint)
  http.ListenAndServe(":8080", nil)
}
