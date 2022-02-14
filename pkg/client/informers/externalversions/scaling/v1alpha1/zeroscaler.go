/*
Copyright 2022 XinYang

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	scalingv1alpha1 "github.com/kzscaler/kzscaler/pkg/apis/scaling/v1alpha1"
	versioned "github.com/kzscaler/kzscaler/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kzscaler/kzscaler/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kzscaler/kzscaler/pkg/client/listers/scaling/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ZeroScalerInformer provides access to a shared informer and lister for
// ZeroScalers.
type ZeroScalerInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ZeroScalerLister
}

type zeroScalerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewZeroScalerInformer constructs a new informer for ZeroScaler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewZeroScalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredZeroScalerInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredZeroScalerInformer constructs a new informer for ZeroScaler type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredZeroScalerInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ScalingV1alpha1().ZeroScalers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ScalingV1alpha1().ZeroScalers(namespace).Watch(context.TODO(), options)
			},
		},
		&scalingv1alpha1.ZeroScaler{},
		resyncPeriod,
		indexers,
	)
}

func (f *zeroScalerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredZeroScalerInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *zeroScalerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&scalingv1alpha1.ZeroScaler{}, f.defaultInformer)
}

func (f *zeroScalerInformer) Lister() v1alpha1.ZeroScalerLister {
	return v1alpha1.NewZeroScalerLister(f.Informer().GetIndexer())
}