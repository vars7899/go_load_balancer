package server_pool

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vars7899/go_load_balancer/server"
	"github.com/vars7899/go_load_balancer/utils"
)

// func AppendServerToPoolTest(t *testing.T) {
// 	sp, _ := server_pool.NewServerPool(utils.RoundRobin)
// 	url, _ := url.Parse("http://localhost:8080")
// 	ns := server.NewServer(url, true, httputil.NewSingleHostReverseProxy(url))
// 	sp.AppendServer(ns)

// 	assert.Equal(t, 1, sp.GetServerPoolSize())
// }

func TestNextValidPeer(t *testing.T) {
	sp, _ := NewServerPool(utils.RoundRobin)
	url, _ := url.Parse("http://localhost:3333")
	b := server.NewServer(url, true, httputil.NewSingleHostReverseProxy(url), 1)
	sp.AppendServer(b)

	url, _ = url.Parse("http://localhost:3334")
	b2 := server.NewServer(url, true, httputil.NewSingleHostReverseProxy(url), 1)
	sp.AppendServer(b2)

	url, _ = url.Parse("http://localhost:3335")
	b3 := server.NewServer(url, true, httputil.NewSingleHostReverseProxy(url), 1)
	sp.AppendServer(b3)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			sp.GetNextValidPeer()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			sp.GetNextValidPeer()
		}
	}()

	wg.Wait()
	assert.Equal(t, b.GetUrlEndpoint().String(), sp.GetNextValidPeer().GetUrlEndpoint().String())
}
