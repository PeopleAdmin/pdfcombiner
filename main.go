// Benchmark parallel downloading and combination of some pdf files.
package main

import (
	"net/http"
	"pdfcombiner/server"
)

func main() {
	server := new(server.CombinerServer)
	http.HandleFunc("/favicon.ico", server.Ping)
	http.HandleFunc("/health_check.html", server.Ping)
	http.HandleFunc("/", server.ProcessJob)
	http.ListenAndServe(":8080", nil)
}
