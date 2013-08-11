// Accept requests to combine PDF documents located on Amazon S3 into
// a single document, upload back to S3, and notify a callback URL.
// TODO fail fast if S3 connection is not usable.
// TODO include timing information in response
// TODO need a way to enable auth in responses
package main

import (
	"flag"
	"launchpad.net/goamz/aws"
	"net/http"
	"pdfcombiner/server"
)

var (
	port       string
	serverMode bool
)

func init() {
	flag.BoolVar(&serverMode, "server", false, "run in server mode")
	flag.StringVar(&port, "port", "8080", "port to listen on")
}

func main() {
	flag.Parse()
	verifyAws()
	if serverMode {
		serve()
	}
}

// Verify AWS credentials are set in the environment:
// - AWS_ACCESS_KEY_ID
// - AWS_SECRET_ACCESS_KEY
func verifyAws() {
	_, err := aws.EnvAuth()
	if err != nil {
		panic(err)
	}
}

// Start a HTTP server listening on `port` to respond
// to JSON-formatted combination requests.
func serve() {
	server := new(server.CombinerServer)
	http.HandleFunc("/health_check", server.Ping)
	http.HandleFunc("/", server.ProcessJob)
	println("Accepting connections on " + port)
	http.ListenAndServe(":"+port, nil)
}
