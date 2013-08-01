// Benchmark parallel downloading and combination of some pdf files.
package main

import (
  "pdfcombiner/combiner"
  "net/http"
  "fmt"
)

func check_params(params map[string] []string ) (docs []string, callback string, ok bool) {
  callback = "http://example.com/postback"
  docs, ok = params["docs"]
  return
}

// Looks for one or more ?docs=FILE params and if found starts combination.
func handle_req(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  params := r.Form
  docs, callback, ok := check_params(params)
  if !ok {
    http.Error(w, "Need some docs and a callback url", http.StatusBadRequest)
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
