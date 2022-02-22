package zeroscaler

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/tracker"

	"github.com/kzscaler/kzscaler/pkg/apis/scaling/v1alpha1"
	"github.com/kzscaler/kzscaler/pkg/autoscaleserver"
	zeroscaledobjreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaledobject"
)

type Reconciler struct {
	kubeClientSet    kubernetes.Interface
	deploymentLister appsv1listers.DeploymentLister
	serviceLister    corev1listers.ServiceLister
	endpointsLister  corev1listers.EndpointsLister

	scaleServer autoscaleserver.ScaleServer
	tracker     tracker.Interface
}

// Check that our Reconciler implements parallelreconciler.Interface
var _ zeroscaledobjreconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, p *v1alpha1.ZeroScaledObject) pkgreconciler.Event {
	logger := logging.FromContext(ctx)

	logger.Infof("reconciling ZeroScaledObject,ns:%s,name:%s", p.Namespace, p.Name)

	svc, err := r.reconcileService(ctx, p)
	if err != nil {
		return fmt.Errorf("reconcile workload err,%s", err)
	}
	svcName := fmt.Sprintf("%s.%s", svc.Name, svc.Namespace)

	replicas, err := r.reconcileWorkload(ctx, p, svcName)
	if err != nil {
		return fmt.Errorf("reconcile workload err,%s", err)
	}

	r.scaleServer.UpdateScaleObj(p)

	p.Status.Replicas = replicas
	p.Status.MarkZeroScaledObjectReady()

	return nil
}

func (r *Reconciler) reconcileWorkload(ctx context.Context, p *v1alpha1.ZeroScaledObject, svcName string) (*int32, error) {
	logger := logging.FromContext(ctx)
	workload := p.Spec.Workload

	err := r.tracker.TrackReference(tracker.Reference{
		APIVersion: workload.APIVersion,
		Kind:       workload.Kind,
		Namespace:  workload.Namespace,
		Name:       workload.Name,
	}, p)

	if err != nil {
		logger.Errorf("track deployment error,%s,%s", workload.Name, err)
	}
	dep, err := r.deploymentLister.Deployments(workload.Namespace).Get(workload.Name)
	if err != nil {
		return nil, fmt.Errorf("no such deployments,%s", err)
	} else {
		// add a scale up function
		r.scaleServer.AddScaleHandler(p, func(i int32) error {

			currrentDep, err := r.deploymentLister.Deployments(workload.Namespace).Get(workload.Name)
			if err != nil {
				return fmt.Errorf("scale up,get dep error:%s", err)
			}
			currrentDep.Spec.Replicas = &i

			_, err = r.kubeClientSet.AppsV1().Deployments(workload.Namespace).Update(ctx, currrentDep, metav1.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("scale up,update dep error:%s", err)
			}
			// wait,until deployments status ok
			toCtx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
			defer cancel()
			_ = wait.PollUntilWithContext(toCtx, 1*time.Second, func(context.Context) (done bool, err error) {
				d, err := r.deploymentLister.Deployments(workload.Namespace).Get(workload.Name)
				if d.Status.ReadyReplicas > 0 {
					return true, nil
				}
				return false, nil
			})

			logger.Infof("scale up,update dep success:%s", err)
			return nil
		})
		return dep.Spec.Replicas, nil
	}
}

func (r *Reconciler) reconcileService(ctx context.Context, p *v1alpha1.ZeroScaledObject) (*v1.Service, error) {
	serviceRef := p.Spec.Service
	svc, err := r.serviceLister.Services(serviceRef.Namespace).Get(serviceRef.Name)
	if err != nil {
		return nil, err
	}
	return svc, nil
}
