package zeroscaler

import (
	"context"

	v1alpha1 "github.com/kzscaler/kzscaler/pkg/apis/scaling/v1alpha1"
	zeroscalerreconciler "github.com/kzscaler/kzscaler/pkg/client/injection/reconciler/scaling/v1alpha1/zeroscaler"
	"k8s.io/client-go/dynamic"
	pkgreconciler "knative.dev/pkg/reconciler"
)

type Reconciler struct {
	// dynamicClientSet allows us to configure pluggable Build objects
	dynamicClientSet dynamic.Interface
}

// Check that our Reconciler implements parallelreconciler.Interface
var _ zeroscalerreconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, p *v1alpha1.ZeroScaler) pkgreconciler.Event {

	return nil
}
