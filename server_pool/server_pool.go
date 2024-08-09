package server_pool

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/pterm/pterm"

	"github.com/vars7899/go_load_balancer/server"
	"github.com/vars7899/go_load_balancer/utils"
)

type ServerPool interface {
	GetServers() []server.Server
	AppendServer(server.Server)
	GetNextValidPeer() server.Server
	GetServerPoolSize() int
	GetTotalPoolWeight() int
}

func NewServerPool(poolStrategy utils.LBStrategy) (ServerPool, error) {
	switch poolStrategy {
	case utils.RoundRobin:
		return &roundRobinPool{
			servers: make([]server.Server, 0),
			current: 0,
		}, nil
	case utils.WeightedRoundRobin:
		return &weightedRoundRobinPool{
			servers: make([]server.Server, 0),
			current: 0,
		}, nil
	case utils.LeastConnection:
		return &LeastConnectionPool{
			servers: make([]server.Server, 0),
			current: 0,
		}, nil
	default:
		return nil, errors.New("pool strategy error: invalid/missing strategy type")
	}
}
func CheckHealth(ctx context.Context, sPool ServerPool) {
	aliveChan := make(chan bool, 1)
	
	for _, s := range sPool.GetServers() {
		s := s
		requestCtx, cancel := context.WithTimeout(ctx, 1000*time.Second)
		defer cancel()

		go server.IsServerAlive(requestCtx, aliveChan, s.GetUrlEndpoint())

		select {
		case <-ctx.Done():
			pterm.DefaultSection.WithTopPadding(0).WithBottomPadding(1).Sprintfln("gracefully closing the url health check")
			return
		case aliveStatus := <-aliveChan:
			s.ModifyAliveFlag(aliveStatus)
		}
	}
	
}

func LaunchServerPoolHealthCheck(ctx context.Context, sp ServerPool) {
	area, _ := pterm.DefaultArea.Start()
	defer area.Stop()

	// areas := make([]pterm.AreaPrinter, len(sp.GetServers()) + 1)
	// for index, _ := range areas {
	// 	area, _ := pterm.DefaultArea.Start()
	// 	areas[index] = *area
	// }
	// defer func(){
	// 	for _, area := range areas{
	// 		area.Stop()
	// 	}
	// }()
	interval := 2
	t := time.NewTicker(time.Duration(interval) * time.Second)
	area.Update(
		pterm.DefaultSection.WithTopPadding(0).WithBottomPadding(1).Sprintfln("Starting Load Balancer Health Monitor..."),
	)
	for {
		select {
		case <-t.C:
			go CheckHealth(ctx, sp)
			go LogHealthDetails(area, sp)
		case <-ctx.Done():
			pterm.DefaultSection.WithTopPadding(0).WithBottomPadding(0).Printf("Gracefully Terminating server pool health check... \n")
			return
		}
	}
}

func LaunchConnectionsGraph(ctx context.Context, sPool ServerPool) {
	t := time.NewTicker(2000 * time.Millisecond)
	defer t.Stop()

	connGroupByServer := make([]float64, 0)
	area, _ := pterm.DefaultArea.WithCenter().Start()
	defer area.Stop()

	for {
		select {
		case <-t.C:
			connGroupByServer = connGroupByServer[:0]

			for _, server := range sPool.GetServers() {
				activeConns := server.GetActiveConnections()
				connGroupByServer = append(connGroupByServer, float64(activeConns))
				// log.Printf("Current Active connections: %d, server: %v", activeConns, server.GetUrlEndpoint())
			}
			ServerInformationCluster(sPool)
		case <-ctx.Done():
			utils.Lg.Info("ending server pool health check...")
			return
		}
	}
}
func ServerInformationCluster(sPool ServerPool) {
	// Define panels in a 2D grid system
	var panels pterm.Panels

	for index, _ := range sPool.GetServers() {
		indexStr := strconv.Itoa(index)

		// Create a new panel and append it to panels
		panel := []pterm.Panel{{Data: indexStr}}
		panels = append(panels, panel)
	}

	// Render the panels with a padding of 5
	_ = pterm.DefaultPanel.WithPanels(panels).WithPadding(5).Render()
}
func LogHealthDetails(area *pterm.AreaPrinter, sPool ServerPool){
	var content string 
	content += pterm.DefaultSection.WithTopPadding(0).WithBottomPadding(1).Sprintfln("Load Balancer Health Monitor")
	for index, s := range sPool.GetServers(){
		status := pterm.DefaultBasicText.Sprint(pterm.Red("DOWN ↓"))
		if s.IsAlive() {
			status = pterm.DefaultBasicText.Sprint(pterm.Green("UP ↑"))
		}
		content += pterm.Sprintfln(
			pterm.Black(pterm.BgLightYellow.Sprintf("  Endpoint (%d) >>  \t", index)) + 
			pterm.Cyan(pterm.Bold.Sprintf("%v\t",s.GetUrlEndpoint())) +
			pterm.White(pterm.Italic.Sprintf("%v Rq\t",s.GetWeight())) +
			pterm.Bold.Sprintf("%v\t",status) +
			pterm.White(pterm.Underscore.Sprintf("%v\t",s.GetActiveConnections())) +
			pterm.White(pterm.Underscore.Sprintf("%v\t",s.GetReqServedCount())),
		)
	}
	content += pterm.Sprintfln("\nRefreshing server health status...")
	area.Update(content)
}

// func LogHealthDetails(areas []pterm.AreaPrinter, sPool ServerPool){
// 	for index, s := range sPool.GetServers(){
// 		status := pterm.DefaultBasicText.Sprint(pterm.Red("DOWN ↓"))
// 		if s.IsAlive() {
// 				status = pterm.DefaultBasicText.Sprint(pterm.Green("UP ↑"))
// 		}
// 		areas[index].Update(
// 			pterm.Sprintfln(
// 				pterm.Black(pterm.BgLightYellow.Sprintf("Endpoint (%d) >>\t", index)) + 
// 				pterm.Cyan(pterm.Bold.Sprintf("%v\t",s.GetUrlEndpoint())) +
// 				pterm.White(pterm.Italic.Sprintf("%v Rq\t",s.GetWeight())) +
// 				pterm.Bold.Sprintf("%v\t",status) +
// 				pterm.White(pterm.Underscore.Sprintf("%v\t",s.GetActiveConnections())) +
// 				pterm.White(pterm.Underscore.Sprintf("%v\t",s.GetReqServedCount())),
// 			),
// 		)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// 	areas[len(areas) - 1].Update(
// 		pterm.Sprintfln("\nRefreshing server health status..."),
// 	)
// }