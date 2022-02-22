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

// FakeZeroScaledObjects implements ZeroScaledObjectInterface
type FakeZeroScaledObjects struct {
	Fake *FakeScalingV1alpha1
	ns   string
}

var zeroscaledobjectsResource = schema.GroupVersionResource{Group: "scaling.xiny.dev", Version: "v1alpha1", Resource: "zeroscaledobjects"}

var zeroscaledobjectsKind = schema.GroupVersionKind{Group: "scaling.xiny.dev", Version: "v1alpha1", Kind: "ZeroScaledObject"}

// Get takes name of the zeroScaledObject, and returns the corresponding zeroScaledObject object, and an error if there is any.
func (c *FakeZeroScaledObjects) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ZeroScaledObject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(zeroscaledobjectsResource, c.ns, name), &v1alpha1.ZeroScaledObject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaledObject), err
}

// List takes label and field selectors, and returns the list of ZeroScaledObjects that match those selectors.
func (c *FakeZeroScaledObjects) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ZeroScaledObjectList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(zeroscaledobjectsResource, zeroscaledobjectsKind, c.ns, opts), &v1alpha1.ZeroScaledObjectList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ZeroScaledObjectList{ListMeta: obj.(*v1alpha1.ZeroScaledObjectList).ListMeta}
	for _, item := range obj.(*v1alpha1.ZeroScaledObjectList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested zeroScaledObjects.
func (c *FakeZeroScaledObjects) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(zeroscaledobjectsResource, c.ns, opts))

}

// Create takes the representation of a zeroScaledObject and creates it.  Returns the server's representation of the zeroScaledObject, and an error, if there is any.
func (c *FakeZeroScaledObjects) Create(ctx context.Context, zeroScaledObject *v1alpha1.ZeroScaledObject, opts v1.CreateOptions) (result *v1alpha1.ZeroScaledObject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(zeroscaledobjectsResource, c.ns, zeroScaledObject), &v1alpha1.ZeroScaledObject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaledObject), err
}

// Update takes the representation of a zeroScaledObject and updates it. Returns the server's representation of the zeroScaledObject, and an error, if there is any.
func (c *FakeZeroScaledObjects) Update(ctx context.Context, zeroScaledObject *v1alpha1.ZeroScaledObject, opts v1.UpdateOptions) (result *v1alpha1.ZeroScaledObject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(zeroscaledobjectsResource, c.ns, zeroScaledObject), &v1alpha1.ZeroScaledObject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaledObject), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeZeroScaledObjects) UpdateStatus(ctx context.Context, zeroScaledObject *v1alpha1.ZeroScaledObject, opts v1.UpdateOptions) (*v1alpha1.ZeroScaledObject, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(zeroscaledobjectsResource, "status", c.ns, zeroScaledObject), &v1alpha1.ZeroScaledObject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaledObject), err
}

// Delete takes name of the zeroScaledObject and deletes it. Returns an error if one occurs.
func (c *FakeZeroScaledObjects) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(zeroscaledobjectsResource, c.ns, name), &v1alpha1.ZeroScaledObject{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeZeroScaledObjects) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(zeroscaledobjectsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ZeroScaledObjectList{})
	return err
}

// Patch applies the patch and returns the patched zeroScaledObject.
func (c *FakeZeroScaledObjects) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ZeroScaledObject, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(zeroscaledobjectsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ZeroScaledObject{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ZeroScaledObject), err
}
