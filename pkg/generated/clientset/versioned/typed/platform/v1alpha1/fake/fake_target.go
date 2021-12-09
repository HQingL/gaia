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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "gaia.io/gaia/pkg/apis/platform/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTargets implements TargetInterface
type FakeTargets struct {
	Fake *FakePlatformV1alpha1
}

var targetsResource = schema.GroupVersionResource{Group: "platform.gaia.io", Version: "v1alpha1", Resource: "targets"}

var targetsKind = schema.GroupVersionKind{Group: "platform.gaia.io", Version: "v1alpha1", Kind: "Target"}

// Get takes name of the target, and returns the corresponding target object, and an error if there is any.
func (c *FakeTargets) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Target, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(targetsResource, name), &v1alpha1.Target{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Target), err
}

// List takes label and field selectors, and returns the list of Targets that match those selectors.
func (c *FakeTargets) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.TargetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(targetsResource, targetsKind, opts), &v1alpha1.TargetList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TargetList{ListMeta: obj.(*v1alpha1.TargetList).ListMeta}
	for _, item := range obj.(*v1alpha1.TargetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested targets.
func (c *FakeTargets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(targetsResource, opts))
}

// Create takes the representation of a target and creates it.  Returns the server's representation of the target, and an error, if there is any.
func (c *FakeTargets) Create(ctx context.Context, target *v1alpha1.Target, opts v1.CreateOptions) (result *v1alpha1.Target, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(targetsResource, target), &v1alpha1.Target{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Target), err
}

// Update takes the representation of a target and updates it. Returns the server's representation of the target, and an error, if there is any.
func (c *FakeTargets) Update(ctx context.Context, target *v1alpha1.Target, opts v1.UpdateOptions) (result *v1alpha1.Target, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(targetsResource, target), &v1alpha1.Target{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Target), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeTargets) UpdateStatus(ctx context.Context, target *v1alpha1.Target, opts v1.UpdateOptions) (*v1alpha1.Target, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(targetsResource, "status", target), &v1alpha1.Target{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Target), err
}

// Delete takes name of the target and deletes it. Returns an error if one occurs.
func (c *FakeTargets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(targetsResource, name), &v1alpha1.Target{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTargets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(targetsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.TargetList{})
	return err
}

// Patch applies the patch and returns the patched target.
func (c *FakeTargets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Target, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(targetsResource, name, pt, data, subresources...), &v1alpha1.Target{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Target), err
}
