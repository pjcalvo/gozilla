package main

import (
	"context"
	"fmt"
	"log"
	"time"

	gozilla "github.com/pjcalvo/gozilla/pkg"
)

func main() {
	// initialize a test suite
	suite := gozilla.NewTestSuite().
		WithDuration(time.Minute * 5).
		WithPlotter().
		WithUsers(2).
		WithThinkTime(time.Second * 1)

	// define test tasks
	tasks := gozilla.Tasks{
		// task 1
		{
			Label: "taks # 1",
			Execute: func(_ context.Context) (any, error) {
				return "1", nil
			},
			ExpectedFunc: func(_ context.Context, something any, err error) error {
				val := something.(string)
				if val != "1" {
					return err
				}
				return nil
			},
		},
		// task 2
		{
			Label: "taks # 2",
			Execute: func(_ context.Context) (any, error) {
				return "error", nil
			},
			ExpectedFunc: func(_ context.Context, something any, err error) error {
				val := something.(string)
				if val == "error" {
					return fmt.Errorf("error happened")
				}
				return nil
			},
		},
	}

	// run the suite
	if err := suite.ExecuteTest(tasks); err != nil {
		log.Fatal(err)
	}
}
