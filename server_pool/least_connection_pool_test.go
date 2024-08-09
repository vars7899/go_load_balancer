package server_pool

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vars7899/go_load_balancer/server"
	"github.com/vars7899/go_load_balancer/utils"
)

func SleepHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
}

var (
	h   = http.HandlerFunc(SleepHandler)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	w   = httptest.NewRecorder()
)

func TestLeastConnectionLB(t *testing.T) {
	dummyServer1 := httptest.NewServer(h)
	defer dummyServer1.Close()
	backend1URL, err := url.Parse(dummyServer1.URL)
	if err != nil {
		t.Fatal(err)
	}

	dummyServer2 := httptest.NewServer(h)
	defer dummyServer2.Close()
	backend2URL, err := url.Parse(dummyServer2.URL)
	if err != nil {
		t.Fatal(err)
	}

	rp1 := httputil.NewSingleHostReverseProxy(backend1URL)
	backend1 := server.NewServer(backend1URL, true ,rp1, 1)

	rp2 := httputil.NewSingleHostReverseProxy(backend2URL)
	backend2 := server.NewServer(backend2URL, true ,rp2, 1)

	serverPool, err := NewServerPool(utils.LeastConnection)
	if err != nil {
		t.Fatal(err)
	}

	serverPool.AppendServer(backend1)
	serverPool.AppendServer(backend2)

	assert.Equal(t, 2, serverPool.GetServerPoolSize())

	var wg sync.WaitGroup
	wg.Add(1)

	peer := serverPool.GetNextValidPeer()
	t.Log(peer.GetUrlEndpoint().String())
	assert.NotNil(t, peer)
	go func() {
		defer wg.Done()
		peer.ServeEndpoint(w, req)
	}()
	time.Sleep(1 * time.Second)
	peer2 := serverPool.GetNextValidPeer()
	t.Log(peer2.GetUrlEndpoint().String())
	connPeer2 := peer2.GetActiveConnections()

	assert.NotNil(t, peer)
	assert.Equal(t, 0, connPeer2)
	assert.NotEqual(t, peer, peer2)

	wg.Wait()
}