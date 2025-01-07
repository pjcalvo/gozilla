package pkg

import (
	"context"
	"time"

	"github.com/pjcalvo/gozilla/internal"
)

type Tasks []Task

type Task struct {
	Execute      ExecuteFunc
	ExpectedFunc ExpectFunc
	Label        string
}

func (t *Task) executeLinearTask(ctx context.Context, threadID int) internal.Result {
	startTime := time.Now()

	result := internal.Result{
		Timestamp: startTime,
		ThreadID:  threadID,
		Label:     t.Label,
	}

	res, err := t.Execute(ctx)
	result.Duration = float64(time.Since(startTime).Seconds())

	// expected function overrides all
	if t.ExpectedFunc != nil {
		if err = t.ExpectedFunc(ctx, res, err); err != nil {
			result.Err = err.Error()
			result.Success = false
			return result
		}
	}

	// this means the error was not handled by the expect function
	if err != nil {
		result.Err = err.Error()
		result.Success = false
		return result
	}

	result.Success = true
	return result
}

// prepareTasks runs the init function
func prepareTasks(tasks []Task) ([]Task, error) {
	// for _, t := range tasks {
	// 	var err error
	// 	t.Request.URL, err = t.Request.URL.Parse(fmt.Sprintf("%s%s", baseURL, t.Request.URL))
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	return tasks, nil
}
