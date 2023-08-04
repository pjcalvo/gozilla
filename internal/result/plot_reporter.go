package result

import (
	"fmt"
	"log"
	"time"

	"github.com/gizak/termui/v3/widgets"

	ui "github.com/gizak/termui/v3"
)

func PlotWorker(results <-chan Result, csvReporterPath string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p4 := widgets.NewParagraph()
	p4.Title = "GoZilla nice reporting"
	p4.Text = "Press [q](fg:blue,mod:bold) to QUIT THE REPORT."
	p4.SetRect(0, 0, 60, 3)
	p4.BorderStyle.Fg = ui.ColorYellow

	p0 := widgets.NewPlot()
	p0.Title = "Results over Time"
	p0.Marker = widgets.MarkerDot
	p0.PlotType = 0
	p0.SetRect(0, 3, 60, 25)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen

	p2 := widgets.NewParagraph()
	p2.Title = "Latest sample: "
	p2.SetRect(0, 25, 60, 28)
	p2.BorderStyle.Fg = ui.ColorYellow

	pLog := widgets.NewParagraph()
	pLog.Title = "Test Results: "
	pLog.Text = fmt.Sprintf("[%s](bg:yellow,mod:bold)", csvReporterPath)
	pLog.SetRect(0, 28, 60, 31)
	pLog.BorderStyle.Fg = ui.ColorCyan

	ui.Render(p4, p0, p2, pLog)

	resultsData := make([]Result, 0)
	uiEvents := ui.PollEvents()
	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent {
				switch e.ID {
				case "q", "<C-c>":
					return
				}
			}
		case data := <-results:
			resultsData = append(resultsData, data)
			xLabels, newData := prepareDataForPlot(resultsData)

			p2.Text = fmt.Sprintf("Response time: [%v](fg:blue,mod:bold), status: [%v](fg:red,mod:bold)", data.Duration, data.Status)

			p0.Data = newData
			p0.DataLabels = xLabels
			ui.Render(p4, p0, p2, pLog)
		}
	}
}

func prepareDataForPlot(results []Result) ([]string, [][]float64) {
	var data [][]float64
	var subdata []float64
	var xLabels []string
	for _, r := range results {
		// xLabels = append(xLabels, "juanito")
		subdata = append(subdata, r.Duration)
	}
	data = append(data, subdata)

	return xLabels, data
}

func parseTime(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		log.Fatalf("failed to parse time: %v", err)
	}
	return t
}
