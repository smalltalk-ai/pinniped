// Copyright 2020 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	authenticationv1alpha1 "go.pinniped.dev/generated/1.19/apis/concierge/authentication/v1alpha1"
	versioned "go.pinniped.dev/generated/1.19/client/clientset/versioned"
	internalinterfaces "go.pinniped.dev/generated/1.19/client/informers/externalversions/internalinterfaces"
	v1alpha1 "go.pinniped.dev/generated/1.19/client/listers/authentication/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// WebhookAuthenticatorInformer provides access to a shared informer and lister for
// WebhookAuthenticators.
type WebhookAuthenticatorInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.WebhookAuthenticatorLister
}

type webhookAuthenticatorInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewWebhookAuthenticatorInformer constructs a new informer for WebhookAuthenticator type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewWebhookAuthenticatorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredWebhookAuthenticatorInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredWebhookAuthenticatorInformer constructs a new informer for WebhookAuthenticator type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredWebhookAuthenticatorInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AuthenticationV1alpha1().WebhookAuthenticators(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.AuthenticationV1alpha1().WebhookAuthenticators(namespace).Watch(context.TODO(), options)
			},
		},
		&authenticationv1alpha1.WebhookAuthenticator{},
		resyncPeriod,
		indexers,
	)
}

func (f *webhookAuthenticatorInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredWebhookAuthenticatorInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *webhookAuthenticatorInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&authenticationv1alpha1.WebhookAuthenticator{}, f.defaultInformer)
}

func (f *webhookAuthenticatorInformer) Lister() v1alpha1.WebhookAuthenticatorLister {
	return v1alpha1.NewWebhookAuthenticatorLister(f.Informer().GetIndexer())
}
