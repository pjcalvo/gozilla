package client

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pjcalvo/gozilla/internal/result"
)

type TestSuite struct {
	duration  time.Duration
	users     int
	baseURL   string
	thinkTime time.Duration
}

func NewTestSuite() *TestSuite {
	return &TestSuite{}
}

func (t *TestSuite) WithDuration(duration time.Duration) *TestSuite {
	t.duration = duration
	return t
}

func (t *TestSuite) WithUsers(users int) *TestSuite {
	t.users = users
	return t
}

func (t *TestSuite) WithBaseURL(baseURL string) *TestSuite {
	t.baseURL = baseURL
	return t
}

func (t *TestSuite) WithThinkTime(thinkTime time.Duration) *TestSuite {
	t.thinkTime = thinkTime
	return t
}

func (t *TestSuite) executeTasksForThread(threadID int, tasks []Task, results chan result.Result) {
	client := http.Client{}
	for {
		select {
		// time out after the test is complete
		case <-time.After(t.duration):
			wg.Done()
			return
		default:
			for _, task := range tasks {
				result := task.executeLinearTask(threadID, client, t.baseURL)
				results <- result
			}
			time.Sleep(t.thinkTime)
		}
	}
}

func (t *TestSuite) ExecuteTest(tasks []Task) error {
	tasks, err := cleanUpTasks(tasks, t.baseURL)
	if err != nil {
		return err
	}

	resultsFile, writer, err := result.GenerateFile()
	if err != nil {
		return err
	}
	results := make(chan result.Result, 5)
	resultsPlotter := make(chan result.Result, 1)

	// hardcoded number of workers to write a file
	for w := 1; w <= 5; w++ {
		go result.WritterWorker(writer, results, resultsPlotter)
	}
	// plotting is done with one
	go result.PlotWorker(resultsPlotter, resultsFile)
	printBoxedMessage(fmt.Sprintf("Starting test... Test Results: %s", resultsFile))

	for i := 0; i < t.users; i++ {
		wg.Add(1)
		go t.executeTasksForThread(i, tasks, results)
	}
	wg.Wait()
	return nil
}
