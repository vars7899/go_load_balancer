package server_pool

import (
	"sync"

	"github.com/vars7899/go_load_balancer/server"
)

type roundRobinPool struct {
	servers []server.Server
	mux     sync.Mutex
	current int
}

func (p *roundRobinPool) GetNextValidPeer() server.Server {
	for i := 0; i < len(p.servers); i++ {
		nextPeer := p.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}
func (p *roundRobinPool) GetServers() []server.Server {
	return p.servers
}
func (p *roundRobinPool) AppendServer(s server.Server) {
	p.servers = append(p.servers, s)
}
func (p *roundRobinPool) GetServerPoolSize() int {
	return len(p.servers)
}
func (p *roundRobinPool) GetTotalPoolWeight() int {
	total := 0
	for _, server := range p.servers {
		total += server.GetWeight()
	}
	return total
}
func (p *roundRobinPool) Rotate() server.Server {
	p.mux.Lock()
	p.current = (p.current + 1) % p.GetServerPoolSize()
	defer p.mux.Unlock()
	return p.servers[p.current]
}
