package scheduler

import (
	"context"
	"net/http"
	"strings"
)

// Scheduler tell kzscaler-proxy what to do
type Scheduler interface {
	Start(ctx context.Context) error
	AddService(name string, replicate int)
	UpdateReplicas(name string, replicate int)
}

type SimpleScheduler struct {
	enabledService map[string]int
}

type SimpleSchedulerArgs struct {
}

func NewScheduler() Scheduler {
	return &SimpleScheduler{
		enabledService: map[string]int{},
	}
}

// Start
// This is a blocking call.
func (s *SimpleScheduler) Start(ctx context.Context) error {
	http.HandleFunc("/enabled", s.getEnabledServiceHandler)
	http.HandleFunc("/zerostate", s.getZeroStateHandler)
	return http.ListenAndServe(":8080", nil)
}

func (s *SimpleScheduler) AddService(name string, replicate int) {
	s.enabledService[name] = replicate

}

func (s *SimpleScheduler) UpdateReplicas(name string, replicate int) {
	s.enabledService[name] = replicate
}

func (s *SimpleScheduler) getZeroStateHandler(w http.ResponseWriter, req *http.Request) {
	services := make([]string, 0)
	for k, v := range s.enabledService {
		if v == 0 {
			services = append(services, k)
		}
	}
	w.WriteHeader(231)
	_, _ = w.Write([]byte(strings.Join(services, "|")))

}
func (s *SimpleScheduler) getEnabledServiceHandler(w http.ResponseWriter, req *http.Request) {
	services := make([]string, 0)
	for k, _ := range s.enabledService {
		services = append(services, k)

	}
	w.WriteHeader(231)
	_, _ = w.Write([]byte(strings.Join(services, "|")))

}
