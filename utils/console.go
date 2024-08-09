package utils

import (
	"fmt"
	"log"
	"strconv"

	"github.com/guptarohit/asciigraph"
	"github.com/jroimartin/gocui"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"go.uber.org/zap"
)

func ConsoleGUI() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		Lg.Fatal("console GUI error", zap.String("msg", err.Error()))
	}
	defer g.Close()

	g.SetManagerFunc(ConsoleLayout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, ConsoleQuit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func ConsoleLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("health", 0, 0, maxX/4, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Heath Status")
	}
	return nil
}

func ConsoleQuit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func ConnectionsGraph(data []float64) {
	graph := asciigraph.Plot(data)

	fmt.Println(graph)
}

// CONSOLE THEME
var Primary = pterm.NewStyle(pterm.Bold, pterm.FgGreen)
var Secondary = pterm.NewStyle(pterm.FgWhite, pterm.Italic)

func PaintInSameArea(fn func()){
	area, _ := pterm.DefaultArea.WithCenter().Start()
	fn()
	area.Stop()
}
func HeaderContent()string{
	headerContent, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("X", pterm.FgMagenta.ToStyle()),
		putils.LettersFromStringWithStyle("-", pterm.FgGray.ToStyle()),
		putils.LettersFromStringWithStyle("BLNC", pterm.FgLightRed.ToStyle())).Srender()
	return headerContent
}

func PaintLoadBalancerConfig(config *Config){
	pterm.DefaultBasicText.Printfln(pterm.LightGreen("STRATEGY: \t\t") + pterm.LightWhite(config.Strategy))
	pterm.DefaultBasicText.Printfln(pterm.LightGreen("CHECK INTERVAL: \t") + pterm.LightWhite(config.HealthCheckInterval))
	pterm.DefaultBasicText.Printfln(pterm.LightGreen("MAX RETRY LIMIT: \t") + pterm.LightWhite(config.MaxRetryLimit))
	pterm.DefaultBasicText.Printfln(pterm.LightGreen("TOTAL SERVERS: \t\t") + pterm.LightWhite(len(config.Servers)))
	pterm.DefaultBasicText.Printfln(pterm.LightGreen("PORT: \t\t\t") + pterm.LightWhite(config.Port))
}

func LoadBalancerConfigContent(config *Config) string{
	var content string 
	content += pterm.DefaultBasicText.Sprintfln(pterm.LightGreen("STRATEGY: \t\t\t") + pterm.LightWhite(config.Strategy))
	content += pterm.DefaultBasicText.Sprintfln(pterm.LightGreen("CHECK INTERVAL: \t\t") + pterm.LightWhite(config.HealthCheckInterval))
	content += pterm.DefaultBasicText.Sprintfln(pterm.LightGreen("MAX RETRY LIMIT: \t\t") + pterm.LightWhite(config.MaxRetryLimit))
	content += pterm.DefaultBasicText.Sprintfln(pterm.LightGreen("TOTAL SERVERS: \t\t") + pterm.LightWhite(len(config.Servers)))
	content += pterm.DefaultBasicText.Sprintfln(pterm.LightGreen("PORT: \t\t\t") + pterm.LightWhite(config.Port))
	return content
}

func ServerDetailsContent(config *Config) string{
	tableData := make(pterm.TableData, len(config.Servers))
	tableData = append(tableData, []string{"Endpoints", "Server Weight"})
	for index := range len(config.Servers) {
		endpoint := []string{
			config.Servers[index],
			strconv.Itoa(config.ServerWeights[index]),
		}
		tableData = append(tableData, endpoint)
	}
	content, _ := pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Srender()
	return content
}
func PaintPanel(config *Config){
	// Define panels in a 2D grid system
	panels := pterm.Panels{
		{
			{Data: ""},
		},
		{
			{Data: HeaderContent()},
			{Data: LoadBalancerConfigContent(config)},
		},
		// {
		// 	{Data: ""},
		// },
		// {
		// 	{Data: ServerDetailsContent(config)},
		// },
		// {
		// 	{Data: 	pterm.DefaultSection.WithTopPadding(0).WithBottomPadding(0).Sprintf("Load Balancer is Now Running.... \n")},
		// },
	}
	_ = pterm.DefaultPanel.WithPanels(panels).WithPadding(5).WithBottomPadding(0).Render()
}