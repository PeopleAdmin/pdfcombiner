// Benchmark parallel downloading and combination of some pdf files.
package main

import (
  "pdfcombiner/server"
  "net/http"
)

func main() {
  http.HandleFunc("/favicon.ico", pdfcombiner.NoopEndpoint)
  http.HandleFunc("/health_check.html", pdfcombiner.NoopEndpoint)
  http.HandleFunc("/", pdfcombiner.CombineEndpoint)
  http.ListenAndServe(":8080", nil)
}
