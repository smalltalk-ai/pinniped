// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "go.pinniped.dev/generated/1.20/apis/concierge/authentication/v1alpha1"
	scheme "go.pinniped.dev/generated/1.20/client/concierge/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WebhookAuthenticatorsGetter has a method to return a WebhookAuthenticatorInterface.
// A group's client should implement this interface.
type WebhookAuthenticatorsGetter interface {
	WebhookAuthenticators() WebhookAuthenticatorInterface
}

// WebhookAuthenticatorInterface has methods to work with WebhookAuthenticator resources.
type WebhookAuthenticatorInterface interface {
	Create(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.CreateOptions) (*v1alpha1.WebhookAuthenticator, error)
	Update(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (*v1alpha1.WebhookAuthenticator, error)
	UpdateStatus(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (*v1alpha1.WebhookAuthenticator, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.WebhookAuthenticator, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.WebhookAuthenticatorList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookAuthenticator, err error)
	WebhookAuthenticatorExpansion
}

// webhookAuthenticators implements WebhookAuthenticatorInterface
type webhookAuthenticators struct {
	client rest.Interface
}

// newWebhookAuthenticators returns a WebhookAuthenticators
func newWebhookAuthenticators(c *AuthenticationV1alpha1Client) *webhookAuthenticators {
	return &webhookAuthenticators{
		client: c.RESTClient(),
	}
}

// Get takes name of the webhookAuthenticator, and returns the corresponding webhookAuthenticator object, and an error if there is any.
func (c *webhookAuthenticators) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	result = &v1alpha1.WebhookAuthenticator{}
	err = c.client.Get().
		Resource("webhookauthenticators").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of WebhookAuthenticators that match those selectors.
func (c *webhookAuthenticators) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WebhookAuthenticatorList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.WebhookAuthenticatorList{}
	err = c.client.Get().
		Resource("webhookauthenticators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested webhookAuthenticators.
func (c *webhookAuthenticators) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("webhookauthenticators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a webhookAuthenticator and creates it.  Returns the server's representation of the webhookAuthenticator, and an error, if there is any.
func (c *webhookAuthenticators) Create(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.CreateOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	result = &v1alpha1.WebhookAuthenticator{}
	err = c.client.Post().
		Resource("webhookauthenticators").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookAuthenticator).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a webhookAuthenticator and updates it. Returns the server's representation of the webhookAuthenticator, and an error, if there is any.
func (c *webhookAuthenticators) Update(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	result = &v1alpha1.WebhookAuthenticator{}
	err = c.client.Put().
		Resource("webhookauthenticators").
		Name(webhookAuthenticator.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookAuthenticator).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *webhookAuthenticators) UpdateStatus(ctx context.Context, webhookAuthenticator *v1alpha1.WebhookAuthenticator, opts v1.UpdateOptions) (result *v1alpha1.WebhookAuthenticator, err error) {
	result = &v1alpha1.WebhookAuthenticator{}
	err = c.client.Put().
		Resource("webhookauthenticators").
		Name(webhookAuthenticator.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookAuthenticator).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the webhookAuthenticator and deletes it. Returns an error if one occurs.
func (c *webhookAuthenticators) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("webhookauthenticators").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *webhookAuthenticators) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("webhookauthenticators").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched webhookAuthenticator.
func (c *webhookAuthenticators) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookAuthenticator, err error) {
	result = &v1alpha1.WebhookAuthenticator{}
	err = c.client.Patch(pt).
		Resource("webhookauthenticators").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
