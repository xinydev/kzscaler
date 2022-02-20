package scheduler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// Scheduler tell kzscaler-proxy what to do
type Scheduler interface {
	Start(ctx context.Context) error
	AddService(name string, replicate int)
	UpdateReplicas(name string, replicate int)
}

type SimpleScheduler struct {
	services map[string]int // zero-scale feature enabled services
	router   *gin.Engine
	logger   *zap.SugaredLogger
}

type SimpleSchedulerArgs struct {
}

func NewScheduler() Scheduler {
	l, _ := zap.NewDevelopment()
	return &SimpleScheduler{
		services: map[string]int{},
		router:   gin.Default(),
		logger:   l.Sugar(),
	}
}

// Start
// This is a blocking call.
func (s *SimpleScheduler) Start(ctx context.Context) error {
	s.router.GET("/service", s.getServicesHandler)
	s.router.GET("/scale_up/:service", s.scaleUpHandler)

	return s.router.Run(":8080")
}

func (s *SimpleScheduler) AddService(name string, replicate int) {
	s.services[name] = replicate

}

func (s *SimpleScheduler) UpdateReplicas(name string, replicate int) {
	s.services[name] = replicate
}

func (s *SimpleScheduler) getServicesHandler(c *gin.Context) {

	services := make([]string, 0)
	for k, v := range s.services {
		services = append(services, fmt.Sprintf("%s%%s", k, v))
	}

	c.String(http.StatusOK, strings.Join(services, "&"))
}

func (s *SimpleScheduler) scaleUpHandler(c *gin.Context) {
	name := c.Param("name")
	s.logger.Infof("service need scale up,%s", name)
	time.Sleep(3 * time.Second)
	c.String(http.StatusOK, "")
}
