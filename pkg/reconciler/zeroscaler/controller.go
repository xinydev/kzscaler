package zeroscaler

import (
	"context"

	"go.uber.org/zap"
	corev1 "k8s.io/api/apps/v1"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/endpoints"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/service"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/tracker"

	zeroscalerinformers "github.com/kzscaler/kzscaler/pkg/client/injection/informers/scaling/v1alpha1/zeroscaler"
	zeroscalerreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaler"
	"github.com/kzscaler/kzscaler/pkg/scheduler"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	deploymentInformer := deployment.Get(ctx)
	serviceInformer := service.Get(ctx)
	endpointsInformer := endpoints.Get(ctx)
	zeroscalerInformer := zeroscalerinformers.Get(ctx)

	r := &Reconciler{
		kubeClientSet:    kubeclient.Get(ctx),
		deploymentLister: deploymentInformer.Lister(),
		serviceLister:    serviceInformer.Lister(),
		endpointsLister:  endpointsInformer.Lister(),
	}
	impl := zeroscalerreconciler.NewImpl(ctx, r)

	s := scheduler.NewScheduler()
	go func() {
		err := s.Start(ctx)
		if err != nil {
			logging.FromContext(ctx).Errorw("Failed starting simple scheduler.", zap.Error(err))
		}
	}()
	r.scheduler = s

	zeroscalerInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Tracker is used to notify us that a ZeroScaler's Deployments has changed so that
	// we can reconcile.
	r.tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	deploymentInformer.Informer().AddEventHandler(controller.HandleAll(
		controller.EnsureTypeMeta(
			r.tracker.OnChanged,
			corev1.SchemeGroupVersion.WithKind("Deployment"))))

	return impl
}
