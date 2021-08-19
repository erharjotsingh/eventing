/*
Copyright 2021 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package containersource

import (
	context "context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	apissourcesv1 "knative.dev/eventing/pkg/apis/sources/v1"
	versioned "knative.dev/eventing/pkg/client/clientset/versioned"
	v1 "knative.dev/eventing/pkg/client/informers/externalversions/sources/v1"
	client "knative.dev/eventing/pkg/client/injection/client"
	factory "knative.dev/eventing/pkg/client/injection/informers/factory"
	sourcesv1 "knative.dev/eventing/pkg/client/listers/sources/v1"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterInformer(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct{}

func withInformer(ctx context.Context) (context.Context, controller.Informer) {
	f := factory.Get(ctx)
	inf := f.Sources().V1().ContainerSources()
	return context.WithValue(ctx, Key{}, inf), inf.Informer()
}

func withDynamicInformer(ctx context.Context) context.Context {
	inf := &wrapper{client: client.Get(ctx)}
	return context.WithValue(ctx, Key{}, inf)
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context) v1.ContainerSourceInformer {
	untyped := ctx.Value(Key{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch knative.dev/eventing/pkg/client/informers/externalversions/sources/v1.ContainerSourceInformer from context.")
	}
	return untyped.(v1.ContainerSourceInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string
}

var _ v1.ContainerSourceInformer = (*wrapper)(nil)
var _ sourcesv1.ContainerSourceLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apissourcesv1.ContainerSource{}, 0, nil)
}

func (w *wrapper) Lister() sourcesv1.ContainerSourceLister {
	return w
}

func (w *wrapper) ContainerSources(namespace string) sourcesv1.ContainerSourceNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apissourcesv1.ContainerSource, err error) {
	lo, err := w.client.SourcesV1().ContainerSources(w.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apissourcesv1.ContainerSource, error) {
	return w.client.SourcesV1().ContainerSources(w.namespace).Get(context.TODO(), name, metav1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
