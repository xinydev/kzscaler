package v1alpha1

import (
	"github.com/kzscaler/kzscaler/pkg/apis/scaling"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: scaling.GroupName, Version: "v1alpha1"}
