/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/suzerain-io/pinniped/generated/1.18/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// WebhookIdentityProviders returns a WebhookIdentityProviderInformer.
	WebhookIdentityProviders() WebhookIdentityProviderInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// WebhookIdentityProviders returns a WebhookIdentityProviderInformer.
func (v *version) WebhookIdentityProviders() WebhookIdentityProviderInformer {
	return &webhookIdentityProviderInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
