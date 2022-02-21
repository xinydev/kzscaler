package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Scheduler tell kzscaler-proxy what to do
type Scheduler interface {
	Start(ctx context.Context) error
	UpdateReplicas(name string, replicate int32)
	AddScaleHandler(name string, h func(int32) error)
	DeleteHandler(name string)
}

type SimpleScheduler struct {
	services     map[string]int32 // zero-scale feature enabled services
	scaleHandler map[string]func(int32) error
	router       *gin.Engine
	logger       *zap.SugaredLogger
	observer     Observer

	svcLock sync.Mutex
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
		observer:     NewPromObserver("http://prometheus.istio-system:9090"),
	}
}

// Start
// This is a blocking call.
func (s *SimpleScheduler) Start(ctx context.Context) error {
	s.router.GET("/service", s.getServicesHandler)
	s.router.GET("/scale_up/:service", s.scaleUpHandler)
	go func() {

		err := wait.PollUntilWithContext(ctx, time.Minute, func(ctx context.Context) (done bool, err error) {
			s.svcLock.Lock()
			defer s.svcLock.Unlock()
			for svc, v := range s.services {
				if v == 0 {
					continue
				}
				name, ns := getServiceNameAndNs(svc)
				cnt, err := s.observer.GetMetrics(name, ns, 3)
				if err != nil {
					s.logger.Errorf("get metrics for %s:%s error:%s", ns, name, err)
				}
				// no traffic for this service,scale to zero
				if cnt < 0.0001 {
					if handler, ok := s.scaleHandler[svc]; ok {
						err := handler(0)
						if err != nil {
							s.logger.Errorf("scale down for %s:%s error:%s", ns, name, err)
						}
					} else {
						s.logger.Errorf("scale down,no such service handler,%s:%s", ns, name)
					}

				}

			}

			return false, nil
		})
		if err != nil {
			s.logger.Errorf("run oberser error:%s", err)
		}

	}()
	return s.router.Run(":8080")
}

func (s *SimpleScheduler) UpdateReplicas(name string, replicate int32) {
	s.svcLock.Lock()
	defer s.svcLock.Unlock()
	s.services[name] = replicate
}

func (s *SimpleScheduler) AddScaleHandler(name string, h func(int32) error) {
	s.scaleHandler[name] = h
}

func (s *SimpleScheduler) DeleteHandler(name string) {
	s.svcLock.Lock()
	defer s.svcLock.Unlock()
	delete(s.scaleHandler, name)
	delete(s.services, name)
}
func (s *SimpleScheduler) getServicesHandler(c *gin.Context) {
	s.svcLock.Lock()
	defer s.svcLock.Unlock()

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

func getServiceNameAndNs(name string) (string, string) {
	names := strings.Split(name, ".")
	return names[0], names[1]
}
