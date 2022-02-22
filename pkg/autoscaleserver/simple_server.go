package autoscaleserver

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

	"github.com/kzscaler/kzscaler/pkg/apis/scaling/v1alpha1"
)

// ScaleServer tell kzscaler-proxy what to do
type ScaleServer interface {
	Start(ctx context.Context) error
	UpdateScaleObj(object *v1alpha1.ZeroScaledObject)
	AddScaleHandler(object *v1alpha1.ZeroScaledObject, h func(int32) error)
	DeleteScaleObj(object *v1alpha1.ZeroScaledObject)
}

type SimpleAutoScaleServer struct {
	objs     map[string]*SimpleScaleObj // zero-scale feature enabled services
	router   *gin.Engine
	logger   *zap.SugaredLogger
	observer Observer

	objLock sync.Mutex
}

type SimpleSchedulerArgs struct {
}

type SimpleScaleObj struct {
	ScaleHandler func(int32) error
	Name         string
	Namespace    string
	Replicas     int32
	StableWindow int
}

func NewAutoScaleServer() ScaleServer {
	l, _ := zap.NewDevelopment()
	return &SimpleAutoScaleServer{
		objs:     map[string]*SimpleScaleObj{},
		router:   gin.Default(),
		logger:   l.Sugar(),
		observer: NewPromObserver("http://prometheus.istio-system:9090"),
	}
}

// Start
// This is a blocking call.
func (s *SimpleAutoScaleServer) Start(ctx context.Context) error {
	s.router.GET("/service", s.getServicesHandler)
	s.router.GET("/scale_up/:service", s.scaleUpHandler)
	go func() {

		err := wait.PollUntilWithContext(ctx, time.Minute, func(ctx context.Context) (done bool, err error) {
			s.objLock.Lock() // TODO(xinydev): lock free
			defer s.objLock.Unlock()
			for name, v := range s.objs {
				if v.Replicas == 0 {
					continue
				}
				interval := max(300, v.StableWindow)

				cnt, err := s.observer.GetMetrics(v.Name, v.Namespace, interval)
				if err != nil {
					s.logger.Errorf("get metrics for %s error:%s", name, err)
				}
				// no traffic for this service,scale to zero
				if cnt < 0.0001 {
					if v.ScaleHandler != nil {
						err := v.ScaleHandler(0)
						if err != nil {
							s.logger.Errorf("scale down for %s error:%s", name, err)
						}
					} else {
						s.logger.Errorf("scale down for %s error:no scale handler", name)

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
func (s *SimpleAutoScaleServer) GetObjName(obj *v1alpha1.ZeroScaledObject) string {
	return fmt.Sprintf("%s.%s", obj.Name, obj.Namespace)
}

func (s *SimpleAutoScaleServer) UpdateScaleObj(obj *v1alpha1.ZeroScaledObject) {
	s.objLock.Lock()
	defer s.objLock.Unlock()
	if v, ok := s.objs[s.GetObjName(obj)]; ok {
		s.updateSimpleObj(obj, v)
	} else {
		s.objs[s.GetObjName(obj)] = s.updateSimpleObj(obj, nil)
	}
}
func (s *SimpleAutoScaleServer) updateSimpleObj(obj *v1alpha1.ZeroScaledObject, simpleObj *SimpleScaleObj) *SimpleScaleObj {

	if simpleObj == nil {
		simpleObj = &SimpleScaleObj{}
	}
	simpleObj.Replicas = *obj.Status.Replicas
	if obj.Spec.Rule != nil && obj.Spec.Rule.StableWindow != nil {
		simpleObj.StableWindow = *obj.Spec.Rule.StableWindow
	}
	return simpleObj

}

func (s *SimpleAutoScaleServer) AddScaleHandler(obj *v1alpha1.ZeroScaledObject, h func(int32) error) {
	s.objLock.Lock()
	defer s.objLock.Unlock()
	if _, ok := s.objs[s.GetObjName(obj)]; !ok {
		s.objs[s.GetObjName(obj)] = s.updateSimpleObj(obj, nil)
	}
	s.objs[s.GetObjName(obj)].ScaleHandler = h
}

func (s *SimpleAutoScaleServer) DeleteScaleObj(obj *v1alpha1.ZeroScaledObject) {
	s.objLock.Lock()
	defer s.objLock.Unlock()
	delete(s.objs, s.GetObjName(obj))
}

// getServicesHandler return serviceA%3&serviceB%0
func (s *SimpleAutoScaleServer) getServicesHandler(c *gin.Context) {
	s.objLock.Lock()
	defer s.objLock.Unlock()

	services := make([]string, 0)
	for k, v := range s.objs {
		services = append(services, fmt.Sprintf("%s%s%d", k, "%", v.Replicas))
	}

	c.String(http.StatusOK, strings.Join(services, "&"))
}

func (s *SimpleAutoScaleServer) scaleUpHandler(c *gin.Context) {
	name := c.Param("service")
	s.logger.Infof("service need scale up,%s", name)

	if scaleObj, ok := s.objs[name]; ok {
		if scaleObj.ScaleHandler != nil {
			err := scaleObj.ScaleHandler(1)
			if err == nil {
				c.String(http.StatusOK, "")
				return
			}
			s.logger.Errorf("scale up error,%s", err)
		}
	}
	c.String(http.StatusServiceUnavailable, "")
}

func getServiceNameAndNs(name string) (string, string) {
	names := strings.Split(name, ".")
	return names[0], names[1]
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
