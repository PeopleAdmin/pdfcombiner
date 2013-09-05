// Package server contains handlers which respond to requests to combine
// pdf files and hands them off to the combiner package.
// TODO more informative validation errors
package server

import (
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/combiner"
	"github.com/PeopleAdmin/pdfcombiner/job"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

// A CombinerServer needs a port to listen on and a WaitGroup to keep
// track of how many background jobs it has spawned that have not yet
// yet completed their work.
type CombinerServer struct {
	port    int
	pending *sync.WaitGroup
}

var invalidMessage = []byte("{\"response\":\"invalid params\"}\n")
var okMessage = []byte("{\"response\":\"ok\"}\n")
var host, _ = os.Hostname()

// Listen starts a HTTP server listening on `Port` to respond to
// JSON-formatted combination requests.
func (c CombinerServer) Listen(listenPort int) {
	c.port = listenPort
	c.pending = new(sync.WaitGroup)
	listener, err := net.Listen("tcp", c.portString())
	if err != nil {
		panic(err)
	}
	println("Accepting connections on " + c.portString())
	c.registerHandlers(listener)
	http.Serve(listener, http.DefaultServeMux)
	println("Waiting for all jobs to finish...")
	c.pending.Wait()
}

// ProcessJob is a handler to recieve a JSON body encoding a Job.  If
// it validates, send it along to be fulfilled, keeping track of the
// in-flight job count.
func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, err := job.NewFromJSON(r.Body)
	logJobReceipt(r, j)
	if err != nil || !j.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(invalidMessage)
		return
	}
	w.Write(okMessage)
	c.registerWorker()
	go func() {
		defer c.unregisterWorker()
		combiner.Combine(j)
	}()
}

// Ping is no-op handler for responding to things like health checks.
// It responds 200 OK with no content to all requests.
func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {
	log.Println(requestInfo(r))
}

func (c CombinerServer) Status(w http.ResponseWriter, r *http.Request) {
	log.Println(requestInfo(r))
	w.Header().Set("Content-Type", "application/json")
	template := `{"host": "%s", "running": %d, "waiting": %d}`+"\n"
	jobs := fmt.Sprintf(template,
		host, combiner.CurrentJobs(), combiner.CurrentWait())
	w.Write([]byte(jobs))
}

func logJobReceipt(r *http.Request, j *job.Job) {
	log.Printf("%v, callback: %v\n", requestInfo(r), j.Callback)
}

func requestInfo(r *http.Request) string {
	return fmt.Sprintf("%v %v from %v", r.Method, r.URL, r.RemoteAddr)
}

// http.ListenAndServe needs a string for the port.
func (c CombinerServer) portString() string {
	return ":" + strconv.Itoa(c.port)
}

// registerWorker increments the count of in-progress jobs
func (c CombinerServer) registerWorker() {
	c.pending.Add(1)
}

// workerDone decrements the count of in-progress jobs.
func (c CombinerServer) unregisterWorker() {
	c.pending.Done()
}

func (c CombinerServer) registerHandlers(listener net.Listener) {
	http.HandleFunc("/health_check", c.Ping)
	http.HandleFunc("/status", c.Status)
	http.HandleFunc("/", c.ProcessJob)
	handleSignals(listener)
}

// On OS shutdown or Ctrl-C, immediately close the tcp listener,
// but let background jobs finish before actually exiting.
func handleSignals(listener net.Listener) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	go func() {
		for _ = range sigs {
			listener.Close()
		}
	}()
}
