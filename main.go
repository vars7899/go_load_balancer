package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/vars7899/go_load_balancer/load_balancer"
	"github.com/vars7899/go_load_balancer/server"
	"github.com/vars7899/go_load_balancer/server_pool"
	"github.com/vars7899/go_load_balancer/utils"
	"go.uber.org/zap"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

func main() {
	// ? <--CONFIGURE-->
	// coreConfig, err := utils.LoadCoreConfig()
	// if err != nil {
		// 	utils.Lg.Fatal(err.Error())
		// }
		
	logger := utils.InitZapLogger()
	defer logger.Sync()
		
	config, err := utils.LoadConfig()
	if err != nil {
		utils.Lg.Fatal(err.Error())
	}
		
	utils.PaintPanel(config)


	// utils.StartMsg(coreConfig)

	// ? <--CLEANUP-->
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ? <-- GENERATE POOL-->
	sPool, err := server_pool.NewServerPool(utils.GetCurrentLBStrategy(config.Strategy))
	if err != nil {
		utils.Lg.Fatal(err.Error())
	}

	// ? <-- INIT LOAD BALANCER-->
	lb := load_balancer.NewLoadBalancer(sPool)

	for index, _url := range config.Servers {
		endpoint, err := url.Parse(_url)
		if err != nil {
			logger.Fatal(err.Error(), zap.String("URL", _url))
		}

		rp := httputil.NewSingleHostReverseProxy(endpoint)
		currentServer := server.NewServer(endpoint, false, rp, config.ServerWeights[index])

		rp.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {

			// logger.Error("error handling the request",
			// 	zap.String("host", endpoint.Host),
			// 	zap.Error(err),
			// )
			pterm.Error.Printfln("Req Error", err)

			currentServer.ModifyAliveFlag(false)

			if !load_balancer.IsRetryAllowed(r) {
				utils.Lg.Info(
					"Max retry attempts reached, terminating",
					zap.String("address", r.RemoteAddr),
					zap.String("path", r.URL.Path),
				)
				http.Error(rw, "Service not available", http.StatusServiceUnavailable)
				return
			}

			pterm.Error.Printfln("Attempting Retry")
			// logger.Info(
			// 	"Attempting retry",
			// 	zap.String("address", r.RemoteAddr),
			// 	zap.String("URL", r.URL.Path),
			// 	zap.Bool("retry", true),
			// )

			lb.ServeEndpoint(
				rw,
				r.WithContext(context.WithValue(
					r.Context(),
					load_balancer.RETRY_ATTEMPTED,
					true),
				),
			)
		}

		sPool.AppendServer(currentServer)
	}

	mux := http.NewServeMux()
	// mux.Handle("/studio/", http.StripPrefix("/studio", http.FileServer(http.Dir("./www/dist"))))
	// mux.HandleFunc("/ws", wsHandler(sPool))
	mux.HandleFunc("/lb", lb.ServeEndpoint)

	baseServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	http.Handle("/log", http.FileServer(http.Dir("./www/dist")))

	// go server_pool.LaunchConnectionsGraph(ctx, sPool)
	// go server_pool.ServerInformationCluster(sPool)

	// utils.PaintInSameArea(
	// 	func(){
	// 		go server_pool.LaunchServerPoolHealthCheck(ctx, sPool)
	// 	},
	// )
	
	go server_pool.LaunchServerPoolHealthCheck(ctx, sPool)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := baseServer.Shutdown(shutdownCtx); err != nil {
			logger.Fatal(err.Error())
		}
		defer cancel()
	}()

	if err := baseServer.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("listen & serve error", zap.Error(err))
	}

}

type ServerData struct {
	URL               string `json:"url"`
	ActiveConnections int    `json:"activeConnections"`
	Weight            int    `json:"weight"`
	ReqServedCount    int    `json:"reqServedCount"`
	Status    		  bool   `json:"status"`
}

// func wsHandler(sPool server_pool.ServerPool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		conn, err := upgrader.Upgrade(w, r, nil)
// 		if err != nil {
// 			utils.Lg.Error("Upgrade error:", zap.Error(err))
// 			return
// 		}
// 		defer conn.Close()

// 		ticker := time.NewTicker(10 * time.Second)
// 		defer ticker.Stop()

// 		for {
// 			select {
// 			case <-ticker.C:
// 				var data []ServerData
// 				for _, server := range sPool.GetServers() {
// 					serverData := ServerData{
// 						URL:               server.GetUrlEndpoint().String(),
// 						ActiveConnections: server.GetActiveConnections(),
// 						Weight:            server.GetWeight(),
// 						ReqServedCount:    server.GetReqServedCount(),
// 					}
// 					data = append(data, serverData)
// 				}
// 				err = conn.WriteJSON(data)
// 				if err != nil {
// 					utils.Lg.Error("Write JSON error:", zap.Error(err))
// 					return
// 				}
// 			case <-r.Context().Done():
// 				return
// 			}
// 		}
// 	}
// }
