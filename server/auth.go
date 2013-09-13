package server

import (
	"bitbucket.org/taruti/http_auth"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const settingsFile = ".pdfcombiner.json"

type configfile struct {
	User string `json:"remote_user"`
	Pass string `json:"remote_password"`
}

var config configfile

func init() {
	if os.Getenv("NO_HTTP_AUTH") != "" {
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
	return config.Pass
})

func authenticate(w http.ResponseWriter, r *http.Request) bool {
	if os.Getenv("NO_HTTP_AUTH") != "" {
		return true
	}
	return auth.Auth(w, r)
}
