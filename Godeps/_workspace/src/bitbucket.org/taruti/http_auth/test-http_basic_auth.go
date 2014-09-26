// Test basic auth

// +build ignore

package main

import (
	"bitbucket.org/taruti/http_auth"
	"net/http"
)

var s = http_auth.NewBasic("MyRealm", func(user string, realm string) string {
	// Replace this with a real lookup function here
	return "mypass"
})

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Auth(w, r) {
		return
	}
	w.Write([]byte("<html><title>Hello</title><h1>Hello</h1></html>"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	// In real life you'd use TLS with http Basic Auth
	http.ListenAndServe(":8080", nil)
}
