package server

import (
	"bitbucket.org/taruti/http_auth"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const settingsFile = ".pdfcombiner.json"

type configfile struct {
	User string `json:"remote_user"`
	Pass string `json:"remote_password"`
}

var config configfile

func init() {
	if os.Getenv("NO_HTTP_AUTH") == "" {
		settingsPath := fmt.Sprintf("%s/%s", os.Getenv("HOME"), settingsFile)
		content, err := os.Open(settingsPath)
		if err != nil {
			panic("Need settings file at " + settingsPath)
		}
		err = json.NewDecoder(content).Decode(&config)
		if err != nil {
			panic("Not JSON:" + settingsPath)
		}
		if config.User == "" || config.Pass == "" {
			panic("User name and password must be configured")
		}
	}
}

var auth = http_auth.NewBasic("Realm", func(user string, realm string) string {
	if user != config.User {
		return randPass()
	}
	return config.Pass
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
