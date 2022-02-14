package zeroscaler

import (
	"context"

	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/clients/dynamicclient"

	zeroscalerreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaler"
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

	return impl
}
