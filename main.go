// Benchmark parallel downloading and combination of some pdf files.
package main

import (
  "pdfcombiner/server"
  "net/http"
)

func main() {
  server := new(pdfcombiner.CombinerServer)
  http.HandleFunc("/favicon.ico", server.Ping)
  http.HandleFunc("/health_check.html", server.Ping)
  http.HandleFunc("/", server.ProcessJob)
  http.ListenAndServe(":8080", nil)
}
