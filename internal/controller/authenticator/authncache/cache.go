// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package authncache implements a cache of active authenticators.
package authncache

import (
	"context"
	"sort"
	"sync"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/klog/v2"

	loginapi "go.pinniped.dev/generated/latest/apis/concierge/login"
	"go.pinniped.dev/internal/constable"
	"go.pinniped.dev/internal/plog"
)

// ErrNoSuchAuthenticator is returned by Cache.AuthenticateTokenCredentialRequest() when the requested authenticator is not configured.
const ErrNoSuchAuthenticator = constable.Error("no such authenticator")

// Cache implements the authenticator.Token interface by multiplexing across a dynamic set of authenticators
// loaded from authenticator resources.
type Cache struct {
	cache sync.Map
}

type Key struct {
	APIGroup string
	Kind     string
	Name     string
}

type Value interface {
	authenticator.Token
}

// New returns an empty cache.
func New() *Cache {
	return &Cache{}
}

// Get an authenticator by key.
func (c *Cache) Get(key Key) Value {
	res, _ := c.cache.Load(key)
	if res == nil {
		return nil
	}
	return res.(Value)
}

// Store an authenticator into the cache.
func (c *Cache) Store(key Key, value Value) {
	c.cache.Store(key, value)
}

// Delete an authenticator from the cache.
func (c *Cache) Delete(key Key) {
	c.cache.Delete(key)
}

// Keys currently stored in the cache.
func (c *Cache) Keys() []Key {
	var result []Key
	c.cache.Range(func(key, _ interface{}) bool {
		result = append(result, key.(Key))
		return true
	})

	// Sort the results for consistency.
	sort.Slice(result, func(i, j int) bool {
		return result[i].APIGroup < result[j].APIGroup ||
			result[i].Kind < result[j].Kind ||
			result[i].Name < result[j].Name
	})
	return result
}

func (c *Cache) AuthenticateTokenCredentialRequest(ctx context.Context, req *loginapi.TokenCredentialRequest) (user.Info, error) {
	// Map the incoming request to a cache key.
	key := Key{
		Name: req.Spec.Authenticator.Name,
		Kind: req.Spec.Authenticator.Kind,
	}
	if req.Spec.Authenticator.APIGroup != nil {
		key.APIGroup = *req.Spec.Authenticator.APIGroup
	}

	val := c.Get(key)
	if val == nil {
		plog.Debug(
			"authenticator does not exist",
			"authenticator", klog.KRef("", key.Name),
			"kind", key.Kind,
			"apiGroup", key.APIGroup,
		)
		return nil, ErrNoSuchAuthenticator
	}

	// The incoming context could have an audience. Since we do not want to handle audiences right now, do not pass it
	// through directly to the authentication webhook.
	ctx = valuelessContext{ctx}

	// Call the selected authenticator.
	resp, authenticated, err := val.AuthenticateToken(ctx, req.Spec.Token)
	if err != nil {
		return nil, err
	}
	if !authenticated {
		return nil, nil
	}

	// Return the user.Info from the response (if it is non-nil).
	var respUser user.Info
	if resp != nil {
		respUser = resp.User
	}
	return respUser, nil
}

type valuelessContext struct{ context.Context }

func (valuelessContext) Value(interface{}) interface{} { return nil }
