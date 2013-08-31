// Accept requests to combine PDF documents located on Amazon S3 into
// a single document, upload back to S3, and notify a callback URL.
// TODO fail fast if S3 connection is not usable.
// TODO need a way to enable auth in responses
package main

import (
	"flag"
	"github.com/PeopleAdmin/pdfcombiner/combiner"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"github.com/PeopleAdmin/pdfcombiner/server"
	"github.com/PeopleAdmin/pdfcombiner/testmode"
	"launchpad.net/goamz/aws"
	"os"
)

var (
	port       int
	serverMode bool
	bucket     string
	daemon     server.CombinerServer
)

func init() {
	if testmode.IsEnabledFast() {
		println("Starting in test mode")
	}
	flag.BoolVar(&serverMode, "server", false, "run in server mode")
	flag.IntVar(&port, "port", 8080, "port to listen on for server mode")
	flag.StringVar(&bucket, "bucket", "", "bucket name to use in standalone mode")
	flag.Parse()
	flag.Usage = func() {
		println("Usage: pdfcombiner [OPTS] [FILE...]")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	verifyAws()
	switch {
	case serverMode:
		daemon.Listen(port)
		println("Shutdown complete")
	default:
		combineSynchronously()
	}
}

// Combine the requested files and return the status to standard out.
func combineSynchronously() {
	checkArgs()
	pdfFiles := flag.Args()
	j, err := job.New(bucket, pdfFiles)
	if err != nil {
		panic(err)
	}
	combiner.Combine(j)
}

func checkArgs() {
	switch {
	case flag.NArg() < 1:
		println("Cannot start in standalone mode with no files to combine.")
		flag.Usage()
	case bucket == "":
		println("Missing argument")
		flag.Usage()
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
