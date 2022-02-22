package autoscaleserver

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Observer interface {
	GetMetrics(name string, ns string, intervalSecond int) (float64, error)
}

type PromObserver struct {
	clientApi v1.API
}

func NewPromObserver(address string) *PromObserver {
	client, _ := api.NewClient(api.Config{
		Address: address,
	})
	return &PromObserver{
		clientApi: v1.NewAPI(client),
	}
}
func (p *PromObserver) GetMetrics(name string, ns string, intervalSecond int) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	q := fmt.Sprintf("delta(istio_requests_total{"+
		"destination_service_name=\"%s\",destination_service_namespace=\"%s\"}[%d])", name, ns, intervalSecond)
	result, _, err := p.clientApi.Query(ctx, q, time.Now())

	if err != nil {
		return -1, err
	}
	if v, ok := result.(model.Vector); ok {
		if v.Len() > 0 {
			return float64(v[0].Value), nil
		}
	}

	return -1, fmt.Errorf("wrong metrics")

}
