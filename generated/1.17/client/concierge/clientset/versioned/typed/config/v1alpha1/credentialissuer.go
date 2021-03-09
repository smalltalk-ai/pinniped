// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"time"

	v1alpha1 "go.pinniped.dev/generated/1.17/apis/concierge/config/v1alpha1"
	scheme "go.pinniped.dev/generated/1.17/client/concierge/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CredentialIssuersGetter has a method to return a CredentialIssuerInterface.
// A group's client should implement this interface.
type CredentialIssuersGetter interface {
	CredentialIssuers() CredentialIssuerInterface
}

// CredentialIssuerInterface has methods to work with CredentialIssuer resources.
type CredentialIssuerInterface interface {
	Create(*v1alpha1.CredentialIssuer) (*v1alpha1.CredentialIssuer, error)
	Update(*v1alpha1.CredentialIssuer) (*v1alpha1.CredentialIssuer, error)
	UpdateStatus(*v1alpha1.CredentialIssuer) (*v1alpha1.CredentialIssuer, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.CredentialIssuer, error)
	List(opts v1.ListOptions) (*v1alpha1.CredentialIssuerList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CredentialIssuer, err error)
	CredentialIssuerExpansion
}

// credentialIssuers implements CredentialIssuerInterface
type credentialIssuers struct {
	client rest.Interface
}

// newCredentialIssuers returns a CredentialIssuers
func newCredentialIssuers(c *ConfigV1alpha1Client) *credentialIssuers {
	return &credentialIssuers{
		client: c.RESTClient(),
	}
}

// Get takes name of the credentialIssuer, and returns the corresponding credentialIssuer object, and an error if there is any.
func (c *credentialIssuers) Get(name string, options v1.GetOptions) (result *v1alpha1.CredentialIssuer, err error) {
	result = &v1alpha1.CredentialIssuer{}
	err = c.client.Get().
		Resource("credentialissuers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CredentialIssuers that match those selectors.
func (c *credentialIssuers) List(opts v1.ListOptions) (result *v1alpha1.CredentialIssuerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.CredentialIssuerList{}
	err = c.client.Get().
		Resource("credentialissuers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested credentialIssuers.
func (c *credentialIssuers) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("credentialissuers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a credentialIssuer and creates it.  Returns the server's representation of the credentialIssuer, and an error, if there is any.
func (c *credentialIssuers) Create(credentialIssuer *v1alpha1.CredentialIssuer) (result *v1alpha1.CredentialIssuer, err error) {
	result = &v1alpha1.CredentialIssuer{}
	err = c.client.Post().
		Resource("credentialissuers").
		Body(credentialIssuer).
		Do().
		Into(result)
	return
}

// Update takes the representation of a credentialIssuer and updates it. Returns the server's representation of the credentialIssuer, and an error, if there is any.
func (c *credentialIssuers) Update(credentialIssuer *v1alpha1.CredentialIssuer) (result *v1alpha1.CredentialIssuer, err error) {
	result = &v1alpha1.CredentialIssuer{}
	err = c.client.Put().
		Resource("credentialissuers").
		Name(credentialIssuer.Name).
		Body(credentialIssuer).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *credentialIssuers) UpdateStatus(credentialIssuer *v1alpha1.CredentialIssuer) (result *v1alpha1.CredentialIssuer, err error) {
	result = &v1alpha1.CredentialIssuer{}
	err = c.client.Put().
		Resource("credentialissuers").
		Name(credentialIssuer.Name).
		SubResource("status").
		Body(credentialIssuer).
		Do().
		Into(result)
	return
}

// Delete takes name of the credentialIssuer and deletes it. Returns an error if one occurs.
func (c *credentialIssuers) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("credentialissuers").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *credentialIssuers) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("credentialissuers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched credentialIssuer.
func (c *credentialIssuers) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CredentialIssuer, err error) {
	result = &v1alpha1.CredentialIssuer{}
	err = c.client.Patch(pt).
		Resource("credentialissuers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
