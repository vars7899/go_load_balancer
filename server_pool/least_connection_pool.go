package server_pool

import (
	"math"
	"sync"

	"github.com/vars7899/go_load_balancer/server"
)

type LeastConnectionPool struct {
	servers 	[]server.Server
	mux 		sync.RWMutex
	current 	int
}
func (p *LeastConnectionPool) GetNextValidPeer() server.Server {
	for i:=0; i<len(p.servers); i++ {
		nextPeer := p.Rotate()
		if nextPeer != nil && nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}
func (p *LeastConnectionPool) GetServers() []server.Server {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.servers
}
func (p *LeastConnectionPool) AppendServer(s server.Server) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.servers = append(p.servers, s)
}
func (p *LeastConnectionPool) GetServerPoolSize() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return len(p.servers)
}
func (p *LeastConnectionPool) GetTotalPoolWeight() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	total := 0
	for _, s := range p.servers {
		total += s.GetWeight()
	}
	return total;
}
func (p *LeastConnectionPool) Rotate() server.Server{
	p.mux.Lock()
	defer p.mux.Unlock()

	leastConnections := math.MaxInt
	serverWithLeastConnection := -1

	for index, s := range p.servers {
		if !s.IsAlive() {
			continue
		}
		activeConn := s.GetActiveConnections()
		if activeConn < leastConnections {
			leastConnections = activeConn
			serverWithLeastConnection = index
		}
	}

	if serverWithLeastConnection != -1{
		p.current = serverWithLeastConnection
		return p.servers[p.current]
	}

	return nil
}