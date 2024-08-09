package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server interface {
	GetActiveConnections() int
	ModifyAliveFlag(bool)
	IsAlive() bool
	GetUrlEndpoint() *url.URL
	ServeEndpoint(http.ResponseWriter, *http.Request)
	GetWeight() int
	GetReqServedCount() int
	ResetReqServedCount()
}

type server struct {
	url               *url.URL
	mux               sync.RWMutex
	activeConnections int
	alive             bool
	reverseProxy      *httputil.ReverseProxy
	weight            int
	reqCount          int
}

func (s *server) GetActiveConnections() int {
	s.mux.RLock()
	ac := s.activeConnections
	defer s.mux.RUnlock()
	return ac
}
func (s *server) ModifyAliveFlag(updatedValue bool) {
	s.mux.Lock()
	s.alive = updatedValue
	s.mux.Unlock()
}
func (s *server) IsAlive() bool {
	s.mux.RLock()
	aliveStatus := s.alive
	defer s.mux.RUnlock()
	return aliveStatus
}
func (s *server) GetUrlEndpoint() *url.URL {
	return s.url
}
func (s *server) GetWeight() int {
	return s.weight
}
func (s *server) ResetReqServedCount() {
	s.mux.Lock()
	s.reqCount = 0
	defer s.mux.Unlock()
}
func (s *server) GetReqServedCount() int {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.reqCount
}
func (s *server) ServeEndpoint(rw http.ResponseWriter, r *http.Request) {
	s.mux.Lock()
	s.activeConnections++
	s.reqCount++
	s.mux.Unlock()

	defer func() {
		s.mux.Lock()
		s.activeConnections--
		s.mux.Unlock()
	}()

	s.reverseProxy.ServeHTTP(rw, r)
}

func NewServer(url *url.URL, alive bool, rp *httputil.ReverseProxy, w int) *server {
	return &server{
		url:          url,
		alive:        alive,
		reverseProxy: rp,
		weight:       w,
	}
}
