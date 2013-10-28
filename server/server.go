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
	"syscall"
	"time"
)

const ShutdownTimeout = 1 * time.Minute

// A CombinerServer needs a port to listen on and a WaitGroup to keep
// track of how many background jobs it has spawned that have not yet
// yet completed their work.
type CombinerServer struct {
	port int
}

var invalidMessage = []byte("{\"response\":\"invalid params\"}\n")
var okMessage = []byte("{\"response\":\"ok\"}\n")
var host, _ = os.Hostname()
var launchTime = time.Now()

// Listen starts a HTTP server listening on `Port` to respond to
// JSON-formatted combination requests.
func (c CombinerServer) Listen(listenPort int) {
	c.port = listenPort
	listener, err := net.Listen("tcp", c.portString())
	if err != nil {
		panic(err)
	}
	log.Println("Accepting connections on " + c.portString())
	c.registerHandlers(listener)
	http.Serve(listener, http.DefaultServeMux)
	log.Println("Waiting for all jobs to finish...")
	waitForJobs()
}

// ProcessJob is a handler to recieve a JSON body encoding a Job.  If
// it validates, send it along to be fulfilled, keeping track of the
// in-flight job count.
func (c CombinerServer) ProcessJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if !authenticate(w, r) {
		return
	}
	j, err := dispatchNewJob(r, c)
	if err != nil || !j.IsValid() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(invalidMessage)
		return
	}
	w.Write(okMessage)
}

// Ping is no-op handler for responding to things like health checks.
// It responds 200 OK with no content to all requests.
func (c CombinerServer) Ping(w http.ResponseWriter, r *http.Request) {
	log.Println(requestInfo(r))
}

func (c CombinerServer) Status(w http.ResponseWriter, r *http.Request) {
	log.Println(requestInfo(r))
	w.Header().Set("Content-Type", "application/json")
	template := `{"host": "%s", "running": %d, "waiting": %d, "total_received": %d, "uptime_minutes": %d }` + "\n"
	jobs := fmt.Sprintf(template,
		host, combiner.CurrentJobs(), combiner.CurrentWait(), combiner.TotalJobsReceived(), uptime())
	w.Write([]byte(jobs))
}

func uptime() int {
	return int(time.Since(launchTime).Minutes())
}

func logJobReceipt(r *http.Request, j *job.Job) {
	log.Printf("%s %v, callback: %s\n", j.Id(), requestInfo(r), j.Callback)
}

func requestInfo(r *http.Request) string {
	return fmt.Sprintf("%v %v from %v", r.Method, r.URL, r.RemoteAddr)
}

// http.ListenAndServe needs a string for the port.
func (c CombinerServer) portString() string {
	return ":" + strconv.Itoa(c.port)
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

func waitForJobs() {
	a := make(chan time.Time)
	go func(b chan<- time.Time) {
		for {
			if combiner.SafeToQuit() {
				b <- time.Now()
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}(a)
	select {
	case <-a:
	case <-time.After(ShutdownTimeout):
	}
}
