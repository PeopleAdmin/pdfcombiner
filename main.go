// Benchmark parallel downloading and combination of some pdf files.
package main

import (
	"fmt"
	"net/http"
	"pdfcombiner/server"
	"runtime"
)

func init() {
	cpus := runtime.NumCPU()
	fmt.Println("init with %s cpus", cpus)
	runtime.GOMAXPROCS(cpus)
}

func main() {
	server := new(server.CombinerServer)
	http.HandleFunc("/favicon.ico", server.Ping)
	http.HandleFunc("/health_check.html", server.Ping)
	http.HandleFunc("/", server.ProcessJob)
	http.ListenAndServe(":8080", nil)
}
