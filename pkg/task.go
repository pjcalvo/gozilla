package client

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pjcalvo/gozilla/internal/result"
)

type Task struct {
	WaitTime     int
	Request      *http.Request
	ExpectedFunc func(*http.Response) error
	Label        string
}

func (t *Task) defineLabel() {
	if t.Label == "" {
		t.Label = fmt.Sprintf("%s - %s", t.Request.Method, t.Request.RequestURI)
	}
}

func (t *Task) executeLinearTask(threadID int, client http.Client, baseURL string) result.Result {
	t.defineLabel()
	startTime := time.Now()

	result := result.Result{
		Timestamp: startTime,
		ThreadID:  threadID,
		Label:     t.Label,
	}

	resp, err := client.Do(t.Request)
	result.Duration = float64(time.Since(startTime).Milliseconds())

	var body []byte

	if err != nil {
		result.Err = err.Error()
		return result
	}

	result.Status = resp.StatusCode

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		result.Err = err.Error()
	}

	result.Bytes = len(body)
	// result.responseMessage = string(body)

	// expected function overrides all
	if t.ExpectedFunc != nil {
		if err = t.ExpectedFunc(resp); err != nil {
			result.Err = err.Error()
			result.Success = false
			return result
		}
	}

	result.Success = true
	return result
}

// cleanUpTasks preappend the baseURL to the tasks
func prepareTasks(tasks []Task, baseURL string) ([]Task, error) {
	for _, t := range tasks {
		var err error
		t.Request.URL, err = t.Request.URL.Parse(fmt.Sprintf("%s%s", baseURL, t.Request.URL))
		if err != nil {
			return nil, err
		}
	}
	return tasks, nil
}
