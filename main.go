// Accept requests to combine PDF documents located on Amazon S3 into
// a single document, upload back to S3, and notify a callback URL.
// TODO fail fast if S3 connection is not usable.
package main

import (
	"launchpad.net/goamz/aws"
	"net/http"
	"pdfcombiner/server"
)

// Verify AWS credentials are set in the environment:
// - AWS_ACCESS_KEY_ID
// - AWS_SECRET_ACCESS_KEY
func init() {
	_, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
}

func main() {
	server := new(server.CombinerServer)
	http.HandleFunc("/favicon.ico", server.Ping)
	http.HandleFunc("/health_check.html", server.Ping)
	http.HandleFunc("/", server.ProcessJob)
	http.ListenAndServe(":8080", nil)
}
