package scheduler

import (
	"context"
	"fmt"
	"net/http"
)

// Scheduler tell kzscaler-proxy what to do
type Scheduler interface {
	Start(ctx context.Context) error
}

type SimpleScheduler struct {
}

type SimpleSchedulerArgs struct {
}

func NewScheduler() Scheduler {
	return &SimpleScheduler{}
}

// Start
// This is a blocking call.
func (s *SimpleScheduler) Start(ctx context.Context) error {
	http.HandleFunc("/", logHandler)
	return http.ListenAndServe(":8080", nil)
}

func logHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("new Request starts")
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Println(w, "%v: %v\n", name, h)
		}
	}
	fmt.Println("new Request ends")
	w.WriteHeader(231)
}
