// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package authncache

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"

	authv1alpha "go.pinniped.dev/generated/latest/apis/concierge/authentication/v1alpha1"
	loginapi "go.pinniped.dev/generated/latest/apis/concierge/login"
	"go.pinniped.dev/internal/mocks/mocktokenauthenticator"
)

func TestCache(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cache := New()
	require.NotNil(t, cache)

	key1 := Key{Name: "authenticator-one"}
	mockToken1 := mocktokenauthenticator.NewMockToken(ctrl)
	cache.Store(key1, mockToken1)
	require.Equal(t, mockToken1, cache.Get(key1))
	require.Equal(t, 1, len(cache.Keys()))

	key2 := Key{Name: "authenticator-two"}
	mockToken2 := mocktokenauthenticator.NewMockToken(ctrl)
	cache.Store(key2, mockToken2)
	require.Equal(t, mockToken2, cache.Get(key2))
	require.Equal(t, 2, len(cache.Keys()))

	for _, key := range cache.Keys() {
		cache.Delete(key)
	}
	require.Zero(t, len(cache.Keys()))

	// Fill the cache back up with a fixed set of keys, but inserted in shuffled order.
	keysInExpectedOrder := []Key{
		{APIGroup: "a", Kind: "a", Name: "a"},
		{APIGroup: "b", Kind: "a", Name: "a"},
		{APIGroup: "b", Kind: "b", Name: "a"},
		{APIGroup: "b", Kind: "b", Name: "b"},
	}
	for tries := 0; tries < 10; tries++ {
		cache := New()
		for _, i := range rand.Perm(len(keysInExpectedOrder)) {
			cache.Store(keysInExpectedOrder[i], nil)
		}

		// Expect that they come back out in sorted order.
		require.Equal(t, keysInExpectedOrder, cache.Keys())
	}
}

func TestAuthenticateTokenCredentialRequest(t *testing.T) {
	t.Parallel()

	validRequest := loginapi.TokenCredentialRequest{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test-namespace",
		},
		Spec: loginapi.TokenCredentialRequestSpec{
			Authenticator: corev1.TypedLocalObjectReference{
				APIGroup: &authv1alpha.SchemeGroupVersion.Group,
				Kind:     "WebhookAuthenticator",
				Name:     "test-name",
			},
			Token: "test-token",
		},
		Status: loginapi.TokenCredentialRequestStatus{},
	}
	validRequestKey := Key{
		APIGroup: *validRequest.Spec.Authenticator.APIGroup,
		Kind:     validRequest.Spec.Authenticator.Kind,
		Name:     validRequest.Spec.Authenticator.Name,
	}

	mockCache := func(t *testing.T, res *authenticator.Response, authenticated bool, err error) *Cache {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		m := mocktokenauthenticator.NewMockToken(ctrl)
		m.EXPECT().AuthenticateToken(audienceFreeContext{}, validRequest.Spec.Token).Return(res, authenticated, err)
		c := New()
		c.Store(validRequestKey, m)
		return c
	}

	t.Run("no such authenticator", func(t *testing.T) {
		c := New()
		res, err := c.AuthenticateTokenCredentialRequest(context.Background(), validRequest.DeepCopy())
		require.EqualError(t, err, "no such authenticator")
		require.Nil(t, res)
	})

	t.Run("authenticator returns error", func(t *testing.T) {
		c := mockCache(t, nil, false, fmt.Errorf("some authenticator error"))
		res, err := c.AuthenticateTokenCredentialRequest(context.Background(), validRequest.DeepCopy())
		require.EqualError(t, err, "some authenticator error")
		require.Nil(t, res)
	})

	t.Run("authenticator returns unauthenticated without error", func(t *testing.T) {
		c := mockCache(t, &authenticator.Response{}, false, nil)
		res, err := c.AuthenticateTokenCredentialRequest(context.Background(), validRequest.DeepCopy())
		require.NoError(t, err)
		require.Nil(t, res)
	})

	t.Run("authenticator returns nil response without error", func(t *testing.T) {
		c := mockCache(t, nil, true, nil)
		res, err := c.AuthenticateTokenCredentialRequest(context.Background(), validRequest.DeepCopy())
		require.NoError(t, err)
		require.Nil(t, res)
	})

	t.Run("authenticator returns response with nil user", func(t *testing.T) {
		c := mockCache(t, &authenticator.Response{}, true, nil)
		res, err := c.AuthenticateTokenCredentialRequest(context.Background(), validRequest.DeepCopy())
		require.NoError(t, err)
		require.Nil(t, res)
	})

	t.Run("context is cancelled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)
		m := mocktokenauthenticator.NewMockToken(ctrl)
		m.EXPECT().AuthenticateToken(gomock.Any(), validRequest.Spec.Token).DoAndReturn(
			func(ctx context.Context, token string) (*authenticator.Response, bool, error) {
				select {
				case <-time.After(2 * time.Second):
					require.Fail(t, "expected to be cancelled")
					return nil, true, nil
				case <-ctx.Done():
					return nil, false, ctx.Err()
				}
			},
		)
		c := New()
		c.Store(validRequestKey, m)

		ctx, cancel := context.WithCancel(context.Background())
		errchan := make(chan error)
		go func() {
			_, err := c.AuthenticateTokenCredentialRequest(ctx, validRequest.DeepCopy())
			errchan <- err
		}()
		cancel()
		require.EqualError(t, <-errchan, "context canceled")
	})

	t.Run("authenticator returns success", func(t *testing.T) {
		userInfo := user.DefaultInfo{
			Name:   "test-user",
			UID:    "test-uid",
			Groups: []string{"test-group-1", "test-group-2"},
			Extra:  map[string][]string{"extra-key-1": {"extra-value-1", "extra-value-2"}},
		}
		c := mockCache(t, &authenticator.Response{User: &userInfo}, true, nil)

		audienceCtx := authenticator.WithAudiences(context.Background(), authenticator.Audiences{"test-audience-1"})
		res, err := c.AuthenticateTokenCredentialRequest(audienceCtx, validRequest.DeepCopy())
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, "test-user", res.GetName())
		require.Equal(t, "test-uid", res.GetUID())
		require.Equal(t, []string{"test-group-1", "test-group-2"}, res.GetGroups())
		require.Equal(t, map[string][]string{"extra-key-1": {"extra-value-1", "extra-value-2"}}, res.GetExtra())
	})
}

type audienceFreeContext struct{}

func (audienceFreeContext) Matches(in interface{}) bool {
	ctx, isCtx := in.(context.Context)
	if !isCtx {
		return false
	}
	_, hasAudiences := authenticator.AudiencesFrom(ctx)
	return !hasAudiences
}

func (audienceFreeContext) String() string {
	return "is a context without authenticator audiences"
}
