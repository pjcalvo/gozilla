package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	wg        = sync.WaitGroup{}
	sleepTime = time.Duration(1 * time.Second)
	testUrl   = "https://gaslicht.com/"
)

type Result struct {
	status          int
	timestamp       time.Time
	err             string
	duration        float64
	label           string
	responseMessage string
	bytes           int
	threadID        int
}

func writterWorker(writer *csv.Writer, results <-chan Result, resultsPlotter chan Result) {
	for r := range results {
		record := []string{
			r.timestamp.Format(time.RFC3339),
			strconv.Itoa(r.threadID),
			strconv.Itoa(r.status),
			strconv.FormatFloat(r.duration, 'f', 2, 64),
			r.label,
			r.responseMessage,
			strconv.Itoa(r.bytes),
			r.err,
		}
		err := writer.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
		resultsPlotter <- r
	}
}

func plotWorker(results <-chan Result) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	p4 := widgets.NewParagraph()
	p4.Title = "Text Box with Wrapping"
	p4.Text = "Press [q](fg:blue,mod:bold) to QUIT THE REPORT."
	p4.SetRect(0, 0, 60, 3)
	p4.BorderStyle.Fg = ui.ColorYellow

	p0 := widgets.NewPlot()
	p0.Title = "Results over Time"
	p0.Marker = widgets.MarkerDot
	p0.PlotType = 0
	p0.DataLabels = []string{"chorizo"}
	p0.SetRect(0, 3, 60, 25)
	p0.AxesColor = ui.ColorWhite
	p0.LineColors[0] = ui.ColorGreen

	p2 := widgets.NewParagraph()
	p2.Title = "Latest sample: "
	p2.SetRect(0, 25, 60, 28)
	p2.BorderStyle.Fg = ui.ColorYellow

	ui.Render(p4, p0, p2)

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

			p2.Text = fmt.Sprintf("Response time: [%v](fg:blue,mod:bold), status: [%v](fg:red,mod:bold)", data.duration, data.status)

			p0.Data = newData
			p0.DataLabels = xLabels
			ui.Render(p4, p0, p2)
		}
	}
}

func executeRequest(threadID int, results chan Result) {
	for {
		result := Result{
			timestamp: time.Now(),
			threadID:  threadID,
		}
		startTime := time.Now()
		resp, err := http.Get(testUrl)
		result.duration = float64(time.Since(startTime).Milliseconds())

		if err != nil {
			result.err = err.Error()
		} else {
			result.status = resp.StatusCode
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				result.err = err.Error()
			} else {
				result.bytes = len(body)
				// result.responseMessage = string(body)
			}
			resp.Body.Close()
		}

		results <- result
		time.Sleep(sleepTime)
	}
}

func generateFile() (string, *csv.Writer, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("result_%s.csv", timestamp)
	// Create a new file
	file, err := os.Create(fileName)
	if err != nil {
		return "", nil, err
	}

	writer := csv.NewWriter(file)

	headers := []string{
		"timestamp",
		"threadID",
		"status",
		"duration (ms)",
		"label",
		"responseMessage",
		"bytes",
		"error",
	}

	err = writer.Write(headers)
	if err != nil {
		return "", nil, err
	}

	writer.Flush()
	return fileName, writer, nil
}

func prepareDataForPlot(results []Result) ([]string, [][]float64) {
	var data [][]float64
	var subdata []float64
	var xLabels []string
	for _, r := range results {
		xLabels = append(xLabels, "juanito")
		subdata = append(subdata, r.duration)
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

func printBoxedMessage(message string) {
	boxWidth := len(message) + 4
	horizontalLine := "+" + repeatChar("-", boxWidth-2) + "+"

	fmt.Println(horizontalLine)
	fmt.Printf("| %s |\n", centerText(message, boxWidth-4))
	fmt.Println(horizontalLine)
}

func repeatChar(char string, times int) string {
	result := ""
	for i := 0; i < times; i++ {
		result += char
	}
	return result
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text
	}

	padding := width - len(text)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding

	return repeatChar(" ", leftPadding) + text + repeatChar(" ", rightPadding)
}

func main() {
	resultsFile, writer, err := generateFile()
	if err != nil {
		log.Fatal(err)
	}
	results := make(chan Result, 5)
	resultsPlotter := make(chan Result, 1)
	for w := 1; w <= 3; w++ {
		go writterWorker(writer, results, resultsPlotter)
	}
	go plotWorker(resultsPlotter)
	printBoxedMessage(fmt.Sprintf("Starting test... Test Results: %s", resultsFile))

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go executeRequest(i, results)
	}
	wg.Wait()
}
