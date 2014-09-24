package http_auth

import "crypto/rand"
import "fmt"
import "sync"
import "time"

type Server struct {
	sync.Mutex
	AuthFun                           func(user, realm string) string
	opaque                            string
	Realm                             string
	ReaperWaitSeconds, ReapTTLSeconds int
	requests                          map[string]*reqState
}
type reqState struct {
	nreq    int
	sectime int64
}

func newServer(realm string, authfun func(user, realm string) string) *Server {
	s := new(Server)
	s.Realm = realm
	s.opaque = newNonce()
	s.AuthFun = authfun
	s.requests = map[string]*reqState{}
	s.ReaperWaitSeconds = 600
	s.ReapTTLSeconds = 3600
	go reaper(s)

	return s
}

func newNonce() string {
	b := make([]byte, 8)
	n, e := rand.Read(b)
	if n != 8 || e != nil {
		panic("rand.Reader failed!")
	}
	return fmt.Sprintf("%x", b)
}

func reaper(s *Server) {
	for {
		time.Sleep(time.Duration(s.ReaperWaitSeconds) * time.Second)
		do_reap(s)
	}
}

func do_reap(s *Server) {
	now := time.Now().UnixNano()

	s.Lock()
	defer s.Unlock()

	for k, v := range s.requests {
		if v.sectime+int64(s.ReapTTLSeconds) < now {
			delete(s.requests, k)
		}
	}
}
