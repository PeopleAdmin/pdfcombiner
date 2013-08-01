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
func handle_req(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  params := r.Form
  docs, callback, ok := check_params(params)
  if !ok {
    http.Error(w, "Need some docs and a callback url", http.StatusBadRequest)
    return
  }
  fmt.Fprintln(w, "Started combination on",docs)
  go pdfcombiner.Combine(docs,callback)
}

// Exists only to prevent calls to favicon when testing the browser
func noopConn(w http.ResponseWriter, r *http.Request) {}

func main() {
  http.HandleFunc("/favicon.ico", noopConn)
  http.HandleFunc("/", handle_req)
  http.ListenAndServe(":8080", nil)
}
