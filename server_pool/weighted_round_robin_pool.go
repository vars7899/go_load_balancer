package server_pool

import (
	"sync"

	"github.com/vars7899/go_load_balancer/server"
)

type weightedRoundRobinPool struct {
	servers []server.Server
	mux     sync.Mutex
	current int
}

func (p *weightedRoundRobinPool) GetNextValidPeer() server.Server {
	for i := 0; i < len(p.servers); i++ {
		nextPeer := p.Rotate()
		// fmt.Println(nextPeer.IsAlive(), nextPeer.GetUrlEndpoint(), nextPeer.GetReqServedCount())
		if nextPeer.IsAlive() {
			return nextPeer
		}
	}
	return nil
}
func (p *weightedRoundRobinPool) GetServers() []server.Server {
	return p.servers
}
func (p *weightedRoundRobinPool) AppendServer(s server.Server) {
	p.servers = append(p.servers, s)
}
func (p *weightedRoundRobinPool) GetServerPoolSize() int {
	return len(p.servers)
}
func (p *weightedRoundRobinPool) GetTotalPoolWeight() int {
	total := 0
	for _, server := range p.servers {
		total += server.GetWeight()
	}
	return total
}
func (p *weightedRoundRobinPool) Rotate() server.Server {
	p.mux.Lock()
	defer p.mux.Unlock()

	currentServer := p.servers[p.current]
	if currentServer.GetReqServedCount() < currentServer.GetWeight() && currentServer.IsAlive() {
		return currentServer
	}

	p.current = (p.current + 1) % p.GetServerPoolSize()
	currentServer.ResetReqServedCount() // reset req count

	return p.servers[p.current]
}
