package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	gozilla "github.com/pjcalvo/gozilla/pkg"
)

func main() {
	// initialize a test suite
	suite := gozilla.NewTestSuite().
		WithBaseURL("https://pjcalvo.github.io").
		WithDuration(time.Minute * 5).
		WithUsers(2).
		WithThinkTime(time.Second * 1)

	// define test tasks
	reqHomePage, _ := http.NewRequest("GET", "/", nil)
	reqAbout, _ := http.NewRequest("GET", "/about", nil)
	tasks := []gozilla.Task{
		{
			Request: reqHomePage,
			Label:   "homepage",
		},
		{
			Request: reqAbout,
			Label:   "about",
			// define expected behavior
			ExpectedFunc: func(r *http.Response) error {
				if r.StatusCode != 200 {
					return errors.New("response code is not expected")
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
