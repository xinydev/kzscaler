package scheduler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

// Scheduler tell kzscaler-proxy what to do
type Scheduler interface {
	Start(ctx context.Context) error
	UpdateReplicas(name string, replicate int32)
	AddScaleHandler(name string, h func(int32) error)
}

type SimpleScheduler struct {
	services     map[string]int32 // zero-scale feature enabled services
	scaleHandler map[string]func(int32) error
	router       *gin.Engine
	logger       *zap.SugaredLogger
}

type SimpleSchedulerArgs struct {
}

func NewScheduler() Scheduler {
	l, _ := zap.NewDevelopment()
	return &SimpleScheduler{
		services:     map[string]int32{},
		scaleHandler: map[string]func(int32) error{},
		router:       gin.Default(),
		logger:       l.Sugar(),
	}
}

// Start
// This is a blocking call.
func (s *SimpleScheduler) Start(ctx context.Context) error {
	s.router.GET("/service", s.getServicesHandler)
	s.router.GET("/scale_up/:service", s.scaleUpHandler)

	return s.router.Run(":8080")
}

func (s *SimpleScheduler) UpdateReplicas(name string, replicate int32) {
	s.services[name] = replicate
}

func (s *SimpleScheduler) AddScaleHandler(name string, h func(int32) error) {
	s.scaleHandler[name] = h
}

func (s *SimpleScheduler) getServicesHandler(c *gin.Context) {

	services := make([]string, 0)
	for k, v := range s.services {
		services = append(services, fmt.Sprintf("%s%s%d", k, "%", v))
	}

	c.String(http.StatusOK, strings.Join(services, "&"))
}

func (s *SimpleScheduler) scaleUpHandler(c *gin.Context) {
	name := c.Param("service")
	s.logger.Infof("service need scale up,%s", name)

	if handler, ok := s.scaleHandler[name]; ok {
		err := handler(1)
		if err != nil {
			s.logger.Errorf("scale up error,%s", err)
			c.String(520, "")
		} else {
			c.String(http.StatusOK, "")
		}
		return
	} else {
		s.logger.Infof("no such service")
	}

	c.String(http.StatusNotFound, "")
}
