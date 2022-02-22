package zeroscaler

import (
	"context"
	"github.com/kzscaler/kzscaler/pkg/autoscaleserver"

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

	zsoinformers "github.com/kzscaler/kzscaler/pkg/client/injection/informers/scaling/v1alpha1/zeroscaledobject"
	zsoreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaledobject"
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
	zsoInformer := zsoinformers.Get(ctx)

	r := &Reconciler{
		kubeClientSet:    kubeclient.Get(ctx),
		deploymentLister: deploymentInformer.Lister(),
		serviceLister:    serviceInformer.Lister(),
		endpointsLister:  endpointsInformer.Lister(),
	}
	impl := zsoreconciler.NewImpl(ctx, r)

	s := autoscaleserver.NewAutoScaleServer()
	go func() {
		err := s.Start(ctx)
		if err != nil {
			logging.FromContext(ctx).Errorw("Failed starting simple scheduler.", zap.Error(err))
		}
	}()
	r.scaleServer = s

	zsoInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Tracker is used to notify us that a ZeroScaler's Deployments has changed so that
	// we can reconcile.
	r.tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	deploymentInformer.Informer().AddEventHandler(controller.HandleAll(
		controller.EnsureTypeMeta(
			r.tracker.OnChanged,
			corev1.SchemeGroupVersion.WithKind("Deployment"))))

	return impl
}
