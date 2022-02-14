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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/kzscaler/kzscaler/pkg/apis/scaling/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeZeroScalers implements ZeroScalerInterface
type FakeZeroScalers struct {
	Fake *FakeScalingV1alpha1
	ns   string
}

var zeroscalersResource = schema.GroupVersionResource{Group: "scaling", Version: "v1alpha1", Resource: "zeroscalers"}

var zeroscalersKind = schema.GroupVersionKind{Group: "scaling", Version: "v1alpha1", Kind: "ZeroScaler"}

// Get takes name of the zeroScaler, and returns the corresponding zeroScaler object, and an error if there is any.
func (c *FakeZeroScalers) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ZeroScaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(zeroscalersResource, c.ns, name), &v1alpha1.ZeroScaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaler), err
}

// List takes label and field selectors, and returns the list of ZeroScalers that match those selectors.
func (c *FakeZeroScalers) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ZeroScalerList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(zeroscalersResource, zeroscalersKind, c.ns, opts), &v1alpha1.ZeroScalerList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ZeroScalerList{ListMeta: obj.(*v1alpha1.ZeroScalerList).ListMeta}
	for _, item := range obj.(*v1alpha1.ZeroScalerList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested zeroScalers.
func (c *FakeZeroScalers) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(zeroscalersResource, c.ns, opts))

}

// Create takes the representation of a zeroScaler and creates it.  Returns the server's representation of the zeroScaler, and an error, if there is any.
func (c *FakeZeroScalers) Create(ctx context.Context, zeroScaler *v1alpha1.ZeroScaler, opts v1.CreateOptions) (result *v1alpha1.ZeroScaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(zeroscalersResource, c.ns, zeroScaler), &v1alpha1.ZeroScaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaler), err
}

// Update takes the representation of a zeroScaler and updates it. Returns the server's representation of the zeroScaler, and an error, if there is any.
func (c *FakeZeroScalers) Update(ctx context.Context, zeroScaler *v1alpha1.ZeroScaler, opts v1.UpdateOptions) (result *v1alpha1.ZeroScaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(zeroscalersResource, c.ns, zeroScaler), &v1alpha1.ZeroScaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaler), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeZeroScalers) UpdateStatus(ctx context.Context, zeroScaler *v1alpha1.ZeroScaler, opts v1.UpdateOptions) (*v1alpha1.ZeroScaler, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(zeroscalersResource, "status", c.ns, zeroScaler), &v1alpha1.ZeroScaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaler), err
}

// Delete takes name of the zeroScaler and deletes it. Returns an error if one occurs.
func (c *FakeZeroScalers) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(zeroscalersResource, c.ns, name, opts), &v1alpha1.ZeroScaler{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeZeroScalers) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(zeroscalersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ZeroScalerList{})
	return err
}

// Patch applies the patch and returns the patched zeroScaler.
func (c *FakeZeroScalers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ZeroScaler, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(zeroscalersResource, c.ns, name, pt, data, subresources...), &v1alpha1.ZeroScaler{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaler), err
}
