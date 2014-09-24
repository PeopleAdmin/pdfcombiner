package http_auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

type BasicServer Server

func NewBasic(realm string, authfun func(user, realm string) string) *BasicServer {
	return (*BasicServer)(newServer(realm, authfun))
}

func encodeB64(msg string) string {
	encoding := base64.StdEncoding
	buf := make([]byte, encoding.EncodedLen(len(msg)))
	encoding.Encode(buf, []byte(msg))
	return fmt.Sprintf("%s", buf)
}

func decodeB64(msg64 string) string {
	encoding := base64.StdEncoding
	buf := make([]byte, encoding.DecodedLen(len(msg64)))
	encoding.Decode(buf, []byte(msg64))
	// Fix bug where buf is padded with zero bytes
	for i, c := range buf {
		//fmt.Printf("DC -> i: %d c: %c (as x: %x, as d:%d) \n", i, c, c, c)
		if int(c) == 0 {
			return fmt.Sprintf("%s", buf[0:i])
		}
	}
	return fmt.Sprintf("%s", buf)
}

func parseBasicAuthHeader(r *http.Request) (string, string) {
	// Check whether we are basic authenticating
	ad := strings.SplitN(strings.Trim(fmt.Sprintf("%s", r.Header["Authorization"]), "[]"), " ", 2)

	if len(ad) != 2 || ad[0] != "Basic" {
		return "", ""
	}

	// Split username, password as kv
	u_p := strings.SplitN(decodeB64(ad[1]), ":", 2)
	if len(u_p) != 2 {
		return "", ""
	}
	return u_p[0], u_p[1]
}

// Authenticates
func (s *BasicServer) Auth(w http.ResponseWriter, r *http.Request) bool {
	var pass string
	username, password := parseBasicAuthHeader(r)
	if username == "" || password == "" {
		goto NeedAuth
	}

	// Call username lookup and test username/password combo
	pass = s.AuthFun(username, s.Realm)

	if pass != password {
		goto NeedAuth
	}

	return true

NeedAuth:
	w.Header().Set("WWW-Authenticate", "Basic")
	w.WriteHeader(401)
	w.Write([]byte(`<!DOCTYPE html><title>Unauthorized</title><h1>401 Unauthorized</h1></html>`))

	return false
}
