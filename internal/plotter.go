package internal

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/gizak/termui/v3/widgets"

	ui "github.com/gizak/termui/v3"
)

func PlotWorker(results <-chan Result, csvReporterPath string, sigChan chan os.Signal) {
	if err := ui.Init(); err != nil {
		slog.Error("failed to initialize termui: %v", err)
		panic(err)
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

	pResults := widgets.NewParagraph()
	pResults.Title = "Test Results: "
	pResults.Text = fmt.Sprintf("[%s](bg:yellow,mod:bold)", csvReporterPath)
	pResults.SetRect(0, 28, 60, 31)
	pResults.BorderStyle.Fg = ui.ColorCyan

	logs := make([][]string, 0)
	logs = append(logs, []string{"timestamp", "entry"})

	pLog := widgets.NewTable()
	pLog.Rows = logs
	pLog.Title = "Logs: "
	pLog.RowSeparator = true
	pLog.BorderStyle = ui.NewStyle(ui.ColorGreen)
	pLog.SetRect(0, 31, 60, 35)
	pLog.FillRow = true

	ui.Render(p4, p0, p2, pResults)

	resultsData := make([]Result, 0)
	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			if e.Type == ui.KeyboardEvent {
				switch e.ID {
				case "q", "<C-c>":
					sigChan <- os.Interrupt
					return
				}
			}

		case data := <-results:
			resultsData = append(resultsData, data)
			xLabels, newData := prepareDataForPlot(resultsData)

			p2.Text = fmt.Sprintf("Response time: [%v](fg:blue,mod:bold), status: [%v](fg:red,mod:bold)", float32(data.Duration), data.Status)

			if data.Err != "" {
				row := []string{
					time.Now().String(), data.Err,
				}
				logs[0] = row
			}
			// slog.Info(strings.Join(logs[1], "as"))

			p0.Data = newData
			p0.DataLabels = xLabels
			pLog.Rows = logs
			ui.Render(p4, p0, p2, pResults, pLog)
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

func parseTime(timestamp string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %v", err)
	}
	return t, nil
}
