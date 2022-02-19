package zeroscaler

import (
	"context"
	"github.com/kzscaler/kzscaler/pkg/scheduler"
	"go.uber.org/zap"
	"knative.dev/pkg/logging"

	zeroscalerreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaler"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	r := &Reconciler{

		dynamicClientSet: dynamicclient.Get(ctx),
	}
	impl := zeroscalerreconciler.NewImpl(ctx, r)

	s := scheduler.NewScheduler()
	go func() {
		err := s.Start(ctx)
		if err != nil {
			logging.FromContext(ctx).Errorw("Failed starting simple scheduler.", zap.Error(err))
		}
	}()
	return impl
}
