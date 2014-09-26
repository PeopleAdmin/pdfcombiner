// Test digest auth

// +build ignore

package main

import (
	"bitbucket.org/taruti/http_auth"
	"net/http"
)

var s = http_auth.NewDigest("MyRealm", func(user, realm string) string {
	return http_auth.CalculateHA1(user, realm, "mypass")
})

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if !s.Auth(w, r) {
		return
	}
	w.Write([]byte("<html><title>Hello</title><h1>Hello</h1></html>"))
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
