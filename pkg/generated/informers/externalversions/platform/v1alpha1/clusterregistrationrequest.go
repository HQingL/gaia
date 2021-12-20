/*
Copyright The Gaia Authors.

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

	platformv1alpha1 "gaia.io/gaia/pkg/apis/platform/v1alpha1"
	versioned "gaia.io/gaia/pkg/generated/clientset/versioned"
	internalinterfaces "gaia.io/gaia/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "gaia.io/gaia/pkg/generated/listers/platform/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ClusterRegistrationRequestInformer provides access to a shared informer and lister for
// ClusterRegistrationRequests.
type ClusterRegistrationRequestInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ClusterRegistrationRequestLister
}

type clusterRegistrationRequestInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewClusterRegistrationRequestInformer constructs a new informer for ClusterRegistrationRequest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterRegistrationRequestInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterRegistrationRequestInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredClusterRegistrationRequestInformer constructs a new informer for ClusterRegistrationRequest type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterRegistrationRequestInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PlatformV1alpha1().ClusterRegistrationRequests().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PlatformV1alpha1().ClusterRegistrationRequests().Watch(context.TODO(), options)
			},
		},
		&platformv1alpha1.ClusterRegistrationRequest{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterRegistrationRequestInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterRegistrationRequestInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterRegistrationRequestInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&platformv1alpha1.ClusterRegistrationRequest{}, f.defaultInformer)
}

func (f *clusterRegistrationRequestInformer) Lister() v1alpha1.ClusterRegistrationRequestLister {
	return v1alpha1.NewClusterRegistrationRequestLister(f.Informer().GetIndexer())
}