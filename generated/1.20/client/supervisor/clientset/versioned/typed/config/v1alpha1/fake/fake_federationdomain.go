// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/1.20/apis/supervisor/config/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFederationDomains implements FederationDomainInterface
type FakeFederationDomains struct {
	Fake *FakeConfigV1alpha1
	ns   string
}

var federationdomainsResource = schema.GroupVersionResource{Group: "config.supervisor.pinniped.dev", Version: "v1alpha1", Resource: "federationdomains"}

var federationdomainsKind = schema.GroupVersionKind{Group: "config.supervisor.pinniped.dev", Version: "v1alpha1", Kind: "FederationDomain"}

// Get takes name of the federationDomain, and returns the corresponding federationDomain object, and an error if there is any.
func (c *FakeFederationDomains) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FederationDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(federationdomainsResource, c.ns, name), &v1alpha1.FederationDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// List takes label and field selectors, and returns the list of FederationDomains that match those selectors.
func (c *FakeFederationDomains) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FederationDomainList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(federationdomainsResource, federationdomainsKind, c.ns, opts), &v1alpha1.FederationDomainList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FederationDomainList{ListMeta: obj.(*v1alpha1.FederationDomainList).ListMeta}
	for _, item := range obj.(*v1alpha1.FederationDomainList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested federationDomains.
func (c *FakeFederationDomains) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(federationdomainsResource, c.ns, opts))

}

// Create takes the representation of a federationDomain and creates it.  Returns the server's representation of the federationDomain, and an error, if there is any.
func (c *FakeFederationDomains) Create(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.CreateOptions) (result *v1alpha1.FederationDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(federationdomainsResource, c.ns, federationDomain), &v1alpha1.FederationDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// Update takes the representation of a federationDomain and updates it. Returns the server's representation of the federationDomain, and an error, if there is any.
func (c *FakeFederationDomains) Update(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.UpdateOptions) (result *v1alpha1.FederationDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(federationdomainsResource, c.ns, federationDomain), &v1alpha1.FederationDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFederationDomains) UpdateStatus(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.UpdateOptions) (*v1alpha1.FederationDomain, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(federationdomainsResource, "status", c.ns, federationDomain), &v1alpha1.FederationDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// Delete takes name of the federationDomain and deletes it. Returns an error if one occurs.
func (c *FakeFederationDomains) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(federationdomainsResource, c.ns, name), &v1alpha1.FederationDomain{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFederationDomains) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(federationdomainsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FederationDomainList{})
	return err
}

// Patch applies the patch and returns the patched federationDomain.
func (c *FakeFederationDomains) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FederationDomain, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(federationdomainsResource, c.ns, name, pt, data, subresources...), &v1alpha1.FederationDomain{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}
