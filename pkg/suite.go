package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pjcalvo/gozilla/internal"
)

type TestSuite struct {
	duration  time.Duration
	users     int
	isPlotter bool
	thinkTime time.Duration
	ctx       context.Context
	cancelCtx func()
}

func NewTestSuite() *TestSuite {
	ctx, cancel := context.WithCancel(context.Background())
	return &TestSuite{
		ctx:       ctx,
		cancelCtx: cancel,
	}
}

func (t *TestSuite) WithDuration(duration time.Duration) *TestSuite {
	t.duration = duration
	t.ctx, t.cancelCtx = context.WithTimeout(t.ctx, duration)
	return t
}

func (t *TestSuite) WithUsers(users int) *TestSuite {
	t.users = users
	return t
}

func (t *TestSuite) WithPlotter() *TestSuite {
	t.isPlotter = true
	return t
}

func (t *TestSuite) WithThinkTime(thinkTime time.Duration) *TestSuite {
	t.thinkTime = thinkTime
	return t
}

func (t *TestSuite) executeTasksForThread(ctx context.Context, threadID int, tasks []Task, results chan internal.Result) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		// time out after the test is complete
		case <-time.After(t.duration):
			t.cancelCtx()
			wg.Done()
			return
		default:
			for _, task := range tasks {
				result := task.executeLinearTask(ctx, threadID)
				results <- result
			}
			time.Sleep(t.thinkTime)
		}
	}
}

func (t *TestSuite) ExecuteTest(tasks []Task) error {
	defer t.cancelCtx()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	tasks, err := prepareTasks(tasks)
	if err != nil {
		return err
	}

	resultsFile, writer, err := internal.GenerateFile()
	if err != nil {
		return err
	}
	results := make(chan internal.Result, 5)
	resultsPlotter := make(chan internal.Result, 1)

	go func() {
		<-sigChan
		t.cancelCtx()

		slog.Info("Received signal, shutting down...")

	}()

	// hardcoded number of workers to write a file
	for w := 1; w <= 5; w++ {
		go internal.WriterWorker(t.ctx, writer, results, resultsPlotter)
	}
	// plotting is done with one
	if t.isPlotter {
		go internal.PlotWorker(t.ctx, resultsPlotter, resultsFile, sigChan)
	}
	printBoxedMessage(fmt.Sprintf("Starting test... Test Results: %s", resultsFile))

	for i := 0; i < t.users; i++ {
		wg.Add(1)
		go t.executeTasksForThread(t.ctx, i, tasks, results)
	}
	wg.Wait()
	return nil
}
