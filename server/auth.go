package server

import (
	"bitbucket.org/taruti/http_auth"
	"fmt"
	"github.com/PeopleAdmin/pdfcombiner/config"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var auth = http_auth.NewBasic("Realm", func(user string, realm string) string {
	if user != config.ListenUser() {
		return randPass()
	}
	return config.ListenPass()
})

func authenticate(w http.ResponseWriter, r *http.Request) (authState bool) {
	if os.Getenv("NO_HTTP_AUTH") != "" {
		return true
	}
	authState = auth.Auth(w, r)
	if !authState {
		log.Println("Authentication failed for ", r)
	}
	return

}

func randPass() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d", rand.Int())
}
