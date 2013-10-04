package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const settingsFile = ".pdfcombiner.json"

type configfile struct {
	// Credentials to secure incoming connections
	ListenUser   string `json:"remote_user"`
	ListenPass   string `json:"remote_password"`
	// Credentials to secure outgoing connections
	TransmitUser string `json:"local_user"`
	TransmitPass string `json:"local_password"`
}

var config configfile

func ListenUser() string {
	return config.ListenUser
}

func ListenPass() string {
	return config.ListenPass
}

func TransmitUser() string {
	return config.TransmitUser
}

func TransmitPass() string {
	return config.TransmitPass
}

func init() {
	if os.Getenv("NO_HTTP_AUTH") == "" {
		settingsPath := fmt.Sprintf("%s/%s", os.Getenv("HOME"), settingsFile)
		content, err := os.Open(settingsPath)
		defer content.Close()
		if err != nil {
			panic("Need settings file at " + settingsPath)
		}
		err = json.NewDecoder(content).Decode(&config)
		if err != nil {
			panic("Not JSON:" + settingsPath)
		}
		if ListenUser() == "" || ListenPass() == "" {
			panic("User name and password must be configured")
		}
		if TransmitUser() == "" || TransmitPass() == "" {
			panic("Transmitting user name and password must be configured")
		}
	}
}
