package load_balancer

import (
	"net/http"

	"github.com/vars7899/go_load_balancer/server_pool"
)

const (
	RETRY_ATTEMPTED = 0
)

func IsRetryAllowed(r *http.Request) bool {
	if _, ok := r.Context().Value(RETRY_ATTEMPTED).(bool); ok {
		return false
	}
	return true
}

type LoadBalancer interface {
	ServeEndpoint(http.ResponseWriter, *http.Request)
}

type loadBalancer struct {
	serverPool server_pool.ServerPool
}

func (lb *loadBalancer) ServeEndpoint(rw http.ResponseWriter, r *http.Request) {
	// timerStart := time.Now()

	availablePeer := lb.serverPool.GetNextValidPeer()
	if availablePeer != nil {
		availablePeer.ServeEndpoint(rw, r)
		// duration := time.Since(timerStart)

		// fmt.Printf("[%v]:[%dms]\n", availablePeer.GetUrlEndpoint().String(), duration)

		return
	}
	http.Error(rw, "lb pool error: service peer not available", http.StatusServiceUnavailable)
}

func NewLoadBalancer(sp server_pool.ServerPool) LoadBalancer {
	return &loadBalancer{
		serverPool: sp,
	}
}
