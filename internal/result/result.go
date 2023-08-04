package result

import "time"

type Result struct {
	Status          int
	Timestamp       time.Time
	Err             string
	Duration        float64
	Label           string
	ResponseMessage string
	Bytes           int
	ThreadID        int
	Success         bool
}
