// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package credentialrequest

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/authentication/user"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/klog/v2"

	loginapi "go.pinniped.dev/generated/latest/apis/concierge/login"
	"go.pinniped.dev/internal/mocks/credentialrequestmocks"
	"go.pinniped.dev/internal/testutil"
)

func TestNew(t *testing.T) {
	r := NewREST(nil, nil, schema.GroupResource{Group: "bears", Resource: "panda"})
	require.NotNil(t, r)
	require.False(t, r.NamespaceScoped())
	require.Equal(t, []string{"pinniped"}, r.Categories())
	require.IsType(t, &loginapi.TokenCredentialRequest{}, r.New())
	require.IsType(t, &loginapi.TokenCredentialRequestList{}, r.NewList())

	ctx := context.Background()

	// check the simple invariants of our no-op list
	list, err := r.List(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.IsType(t, &loginapi.TokenCredentialRequestList{}, list)
	require.Equal(t, "0", list.(*loginapi.TokenCredentialRequestList).ResourceVersion)
	require.NotNil(t, list.(*loginapi.TokenCredentialRequestList).Items)
	require.Len(t, list.(*loginapi.TokenCredentialRequestList).Items, 0)

	// make sure we can turn lists into tables if needed
	table, err := r.ConvertToTable(ctx, list, nil)
	require.NoError(t, err)
	require.NotNil(t, table)
	require.Equal(t, "0", table.ResourceVersion)
	require.Nil(t, table.Rows)

	// exercise group resource - force error by passing a runtime.Object that does not have an embedded object meta
	_, err = r.ConvertToTable(ctx, &metav1.APIGroup{}, nil)
	require.Error(t, err, "the resource panda.bears does not support being converted to a Table")
}

func TestCreate(t *testing.T) {
	spec.Run(t, "create", func(t *testing.T, when spec.G, it spec.S) {
		var r *require.Assertions
		var ctrl *gomock.Controller
		var logger *testutil.TranscriptLogger

		it.Before(func() {
			r = require.New(t)
			ctrl = gomock.NewController(t)
			logger = testutil.NewTranscriptLogger(t)
			klog.SetLogger(logger) // this is unfortunately a global logger, so can't run these tests in parallel :(
		})

		it.After(func() {
			klog.SetLogger(nil)
			ctrl.Finish()
		})

		it("CreateSucceedsWhenGivenATokenAndTheWebhookAuthenticatesTheToken", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req).
				Return(&user.DefaultInfo{
					Name:   "test-user",
					UID:    "test-user-uid",
					Groups: []string{"test-group-1", "test-group-2"},
				}, nil)

			issuer := credentialrequestmocks.NewMockCertIssuer(ctrl)
			issuer.EXPECT().IssuePEM(
				pkix.Name{
					CommonName:   "test-user",
					Organization: []string{"test-group-1", "test-group-2"}},
				[]string{},
				5*time.Minute,
			).Return([]byte("test-cert"), []byte("test-key"), nil)

			storage := NewREST(requestAuthenticator, issuer, schema.GroupResource{})

			response, err := callCreate(context.Background(), storage, req)

			r.NoError(err)
			r.IsType(&loginapi.TokenCredentialRequest{}, response)

			expires := response.(*loginapi.TokenCredentialRequest).Status.Credential.ExpirationTimestamp
			r.NotNil(expires)
			r.InDelta(time.Now().Add(5*time.Minute).Unix(), expires.Unix(), 5)
			response.(*loginapi.TokenCredentialRequest).Status.Credential.ExpirationTimestamp = metav1.Time{}

			r.Equal(response, &loginapi.TokenCredentialRequest{
				Status: loginapi.TokenCredentialRequestStatus{
					Credential: &loginapi.ClusterCredential{
						ExpirationTimestamp:   metav1.Time{},
						ClientCertificateData: "test-cert",
						ClientKeyData:         "test-key",
					},
				},
			})
			requireOneLogStatement(r, logger, `"success" userID:test-user-uid,authenticated:true`)
		})

		it("CreateFailsWithValidTokenWhenCertIssuerFails", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req).
				Return(&user.DefaultInfo{
					Name:   "test-user",
					Groups: []string{"test-group-1", "test-group-2"},
				}, nil)

			issuer := credentialrequestmocks.NewMockCertIssuer(ctrl)
			issuer.EXPECT().
				IssuePEM(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, nil, fmt.Errorf("some certificate authority error"))

			storage := NewREST(requestAuthenticator, issuer, schema.GroupResource{})

			response, err := callCreate(context.Background(), storage, req)
			requireSuccessfulResponseWithAuthenticationFailureMessage(t, err, response)
			requireOneLogStatement(r, logger, `"failure" failureType:cert issuer,msg:some certificate authority error`)
		})

		it("CreateSucceedsWithAnUnauthenticatedStatusWhenGivenATokenAndTheWebhookReturnsNilUser", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req).Return(nil, nil)

			storage := NewREST(requestAuthenticator, nil, schema.GroupResource{})

			response, err := callCreate(context.Background(), storage, req)

			requireSuccessfulResponseWithAuthenticationFailureMessage(t, err, response)
			requireOneLogStatement(r, logger, `"success" userID:<none>,authenticated:false`)
		})

		it("CreateSucceedsWithAnUnauthenticatedStatusWhenWebhookFails", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req).
				Return(nil, errors.New("some webhook error"))

			storage := NewREST(requestAuthenticator, nil, schema.GroupResource{})

			response, err := callCreate(context.Background(), storage, req)

			requireSuccessfulResponseWithAuthenticationFailureMessage(t, err, response)
			requireOneLogStatement(r, logger, `"failure" failureType:token authentication,msg:some webhook error`)
		})

		it("CreateSucceedsWithAnUnauthenticatedStatusWhenWebhookReturnsAnEmptyUsername", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req).
				Return(&user.DefaultInfo{Name: ""}, nil)

			storage := NewREST(requestAuthenticator, nil, schema.GroupResource{})

			response, err := callCreate(context.Background(), storage, req)

			requireSuccessfulResponseWithAuthenticationFailureMessage(t, err, response)
			requireOneLogStatement(r, logger, `"success" userID:,authenticated:false`)
		})

		it("CreateFailsWhenGivenTheWrongInputType", func() {
			notACredentialRequest := runtime.Unknown{}
			response, err := NewREST(nil, nil, schema.GroupResource{}).Create(
				genericapirequest.NewContext(),
				&notACredentialRequest,
				rest.ValidateAllObjectFunc,
				&metav1.CreateOptions{})

			requireAPIError(t, response, err, apierrors.IsBadRequest, "not a TokenCredentialRequest")
			requireOneLogStatement(r, logger, `"failure" failureType:request validation,msg:not a TokenCredentialRequest`)
		})

		it("CreateFailsWhenTokenValueIsEmptyInRequest", func() {
			storage := NewREST(nil, nil, schema.GroupResource{})
			response, err := callCreate(context.Background(), storage, credentialRequest(loginapi.TokenCredentialRequestSpec{
				Token: "",
			}))

			requireAPIError(t, response, err, apierrors.IsInvalid,
				`.pinniped.dev "request name" is invalid: spec.token.value: Required value: token must be supplied`)
			requireOneLogStatement(r, logger, `"failure" failureType:request validation,msg:token must be supplied`)
		})

		it("CreateFailsWhenValidationFails", func() {
			storage := NewREST(nil, nil, schema.GroupResource{})
			response, err := storage.Create(
				context.Background(),
				validCredentialRequest(),
				func(ctx context.Context, obj runtime.Object) error {
					return fmt.Errorf("some validation error")
				},
				&metav1.CreateOptions{})
			r.Nil(response)
			r.EqualError(err, "some validation error")
			requireOneLogStatement(r, logger, `"failure" failureType:validation webhook,msg:some validation error`)
		})

		it("CreateDoesNotAllowValidationFunctionToMutateRequest", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req.DeepCopy()).
				Return(&user.DefaultInfo{Name: "test-user"}, nil)

			storage := NewREST(requestAuthenticator, successfulIssuer(ctrl), schema.GroupResource{})
			response, err := storage.Create(
				context.Background(),
				req,
				func(ctx context.Context, obj runtime.Object) error {
					credentialRequest, _ := obj.(*loginapi.TokenCredentialRequest)
					credentialRequest.Spec.Token = "foobaz"
					return nil
				},
				&metav1.CreateOptions{})
			r.NoError(err)
			r.NotEmpty(response)
		})

		it("CreateDoesNotAllowValidationFunctionToSeeTheActualRequestToken", func() {
			req := validCredentialRequest()

			requestAuthenticator := credentialrequestmocks.NewMockTokenCredentialRequestAuthenticator(ctrl)
			requestAuthenticator.EXPECT().AuthenticateTokenCredentialRequest(gomock.Any(), req.DeepCopy()).
				Return(&user.DefaultInfo{Name: "test-user"}, nil)

			storage := NewREST(requestAuthenticator, successfulIssuer(ctrl), schema.GroupResource{})
			validationFunctionWasCalled := false
			var validationFunctionSawTokenValue string
			response, err := storage.Create(
				context.Background(),
				req,
				func(ctx context.Context, obj runtime.Object) error {
					credentialRequest, _ := obj.(*loginapi.TokenCredentialRequest)
					validationFunctionWasCalled = true
					validationFunctionSawTokenValue = credentialRequest.Spec.Token
					return nil
				},
				&metav1.CreateOptions{})
			r.NoError(err)
			r.NotEmpty(response)
			r.True(validationFunctionWasCalled)
			r.Empty(validationFunctionSawTokenValue)
		})

		it("CreateFailsWhenRequestOptionsDryRunIsNotEmpty", func() {
			response, err := NewREST(nil, nil, schema.GroupResource{}).Create(
				genericapirequest.NewContext(),
				validCredentialRequest(),
				rest.ValidateAllObjectFunc,
				&metav1.CreateOptions{
					DryRun: []string{"some dry run flag"},
				})

			requireAPIError(t, response, err, apierrors.IsInvalid,
				`.pinniped.dev "request name" is invalid: dryRun: Unsupported value: []string{"some dry run flag"}`)
			requireOneLogStatement(r, logger, `"failure" failureType:request validation,msg:dryRun not supported`)
		})

		it("CreateFailsWhenNamespaceIsNotEmpty", func() {
			response, err := NewREST(nil, nil, schema.GroupResource{}).Create(
				genericapirequest.WithNamespace(genericapirequest.NewContext(), "some-ns"),
				validCredentialRequest(),
				rest.ValidateAllObjectFunc,
				&metav1.CreateOptions{})

			requireAPIError(t, response, err, apierrors.IsBadRequest, `namespace is not allowed on TokenCredentialRequest: some-ns`)
			requireOneLogStatement(r, logger, `"failure" failureType:request validation,msg:namespace is not allowed`)
		})
	}, spec.Sequential())
}

func requireOneLogStatement(r *require.Assertions, logger *testutil.TranscriptLogger, messageContains string) {
	transcript := logger.Transcript()
	r.Len(transcript, 1)
	r.Equal("info", transcript[0].Level)
	r.Contains(transcript[0].Message, messageContains)
}

func callCreate(ctx context.Context, storage *REST, obj runtime.Object) (runtime.Object, error) {
	return storage.Create(
		ctx,
		obj,
		rest.ValidateAllObjectFunc,
		&metav1.CreateOptions{
			DryRun: []string{},
		})
}

func validCredentialRequest() *loginapi.TokenCredentialRequest {
	return validCredentialRequestWithToken("some token")
}

func validCredentialRequestWithToken(token string) *loginapi.TokenCredentialRequest {
	return credentialRequest(loginapi.TokenCredentialRequestSpec{Token: token})
}

func credentialRequest(spec loginapi.TokenCredentialRequestSpec) *loginapi.TokenCredentialRequest {
	return &loginapi.TokenCredentialRequest{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: "request name",
		},
		Spec: spec,
	}
}

func requireAPIError(t *testing.T, response runtime.Object, err error, expectedErrorTypeChecker func(err error) bool, expectedErrorMessage string) {
	t.Helper()
	require.Nil(t, response)
	require.True(t, expectedErrorTypeChecker(err))
	var status apierrors.APIStatus
	errors.As(err, &status)
	require.Contains(t, status.Status().Message, expectedErrorMessage)
}

func requireSuccessfulResponseWithAuthenticationFailureMessage(t *testing.T, err error, response runtime.Object) {
	t.Helper()
	require.NoError(t, err)
	require.Equal(t, response, &loginapi.TokenCredentialRequest{
		Status: loginapi.TokenCredentialRequestStatus{
			Credential: nil,
			Message:    stringPtr("authentication failed"),
		},
	})
}

func successfulIssuer(ctrl *gomock.Controller) CertIssuer {
	issuer := credentialrequestmocks.NewMockCertIssuer(ctrl)
	issuer.EXPECT().
		IssuePEM(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]byte("test-cert"), []byte("test-key"), nil)
	return issuer
}

func stringPtr(s string) *string {
	return &s
}
