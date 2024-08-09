package server

import (
	"context"
	"net"
	"net/url"

	"github.com/vars7899/go_load_balancer/utils"
	"go.uber.org/zap"
)

func IsServerAlive(ctx context.Context, aliveChan chan bool, u *url.URL) {
	var d net.Dialer

	con, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		// utils.Lg.Fatal(err.Error())
		aliveChan <- false
		return
	}

	defer func() {
		err = con.Close()
		if err != nil {
			utils.Lg.Fatal("connection not closed: ", zap.String("msg", err.Error()))
		}
	}()

	aliveChan <- true
}
