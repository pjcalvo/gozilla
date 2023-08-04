package result

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func GenerateFile() (string, *csv.Writer, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("%s_results.csv", timestamp)
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
		"success",
		"error",
	}

	err = writer.Write(headers)
	if err != nil {
		return "", nil, err
	}

	writer.Flush()
	return fileName, writer, nil
}

func WritterWorker(writer *csv.Writer, results <-chan Result, resultsPlotter chan Result) {
	for r := range results {
		record := []string{
			r.Timestamp.Format(time.RFC3339),
			strconv.Itoa(r.ThreadID),
			strconv.Itoa(r.Status),
			strconv.FormatFloat(r.Duration, 'f', 2, 64),
			r.Label,
			r.ResponseMessage,
			strconv.Itoa(r.Bytes),
			strconv.FormatBool(r.Success),
			r.Err,
		}
		err := writer.Write(record)
		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()
		resultsPlotter <- r
	}
}
