package http_auth

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DigestServer Server

func NewDigest(realm string, authfun func(user, realm string) string) *DigestServer {
	return (*DigestServer)(newServer(realm, authfun))
}

func CalculateHA1(user, realm, pass string) string {
	h := md5.New()

	fmt.Fprintf(h, "%s:%s:%s", user, realm, pass)
	ha1 := h.Sum(nil)

	return fmt.Sprintf("%x", ha1)
}
func calculateResponse(ha1, method, uri, nonce, cnonce string, nrequest int) string {
	h := md5.New()

	fmt.Fprintf(h, "%s:%s", method, uri)
	ha2 := h.Sum(nil)
	h.Reset()

	fmt.Fprintf(h, "%s:%s:%08x:%s:auth:%x", ha1, nonce, nrequest, cnonce, ha2)
	r := h.Sum(nil)

	return fmt.Sprintf("%x", r)
}

func parseAuthHeader(r *http.Request) map[string]string {
	// Check whether we are digest authenticating
	ad := strings.SplitN(strings.Trim(fmt.Sprintf("%s", r.Header["Authorization"]), "[]"), " ", 2)
	if len(ad) != 2 || ad[0] != "Digest" {
		return nil
	}
	// Split values
	kv := map[string]string{}
	for _, kvr := range strings.Split(ad[1], ",") {
		pair := strings.SplitN(kvr, "=", 2)
		if len(pair) != 2 {
			return nil
		}
		kv[strings.Trim(pair[0], " \"")] = strings.Trim(pair[1], " \"")
	}
	return kv
}

// Authenticates
func (s *DigestServer) Auth(w http.ResponseWriter, r *http.Request) bool {
	var additionalResp string
	var nreq uint64
	var err error
	var response string
	var now int64
	var ha1 string
	var rs *reqState
	var found bool
	kv := parseAuthHeader(r)
	if kv == nil {
		goto NeedAuth
	}
	// Check opaque
	if s.opaque != kv["opaque"] {
		goto NeedAuth
	}
	// Get nreq
	nreq, err = strconv.ParseUint(kv["nc"], 16, 64)
	if err != nil {
		goto NeedAuth
	}
	// Calculate r
	ha1 = s.AuthFun(kv["username"], s.Realm)
	response = calculateResponse(ha1, r.Method, r.RequestURI, kv["nonce"], kv["cnonce"], int(nreq))

	if response != kv["response"] {
		goto NeedAuth
	}

	now = time.Now().UnixNano()

	s.Lock()
	defer s.Unlock()

	// Check nreq
	rs, found = s.requests[kv["nonce"]]
	if !found {
		goto NeedAuth
	}
	if rs.nreq < int(nreq) {
		additionalResp = ` stale="true"`
		goto NeedAuth
	}
	rs.nreq++
	rs.sectime = now

	return true

NeedAuth:
	nonce := newNonce()
	s.requests[nonce] = &reqState{1, time.Now().UnixNano()}
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Digest realm="%s", qop="auth", nonce="%s", opaque="%s"%s`, s.Realm, nonce, s.opaque, additionalResp))
	w.WriteHeader(401)
	w.Write([]byte(`<!DOCTYPE html><title>Unauthorized</title><h1>401 Unauthorized</h1></html>`))

	return false
}

func (s *Server) Logout(r *http.Request) {
	kv := parseAuthHeader(r)
	s.Lock()
	delete(s.requests, kv["nonce"])
	s.Unlock()
}
