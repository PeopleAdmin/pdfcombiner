// Accept requests to combine PDF documents located on Amazon S3 into
// a single document, upload back to S3, and notify a callback URL.
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
