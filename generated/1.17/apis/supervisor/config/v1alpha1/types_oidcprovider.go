// Copyright 2020 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:Enum=Success;Duplicate;Invalid
type OIDCProviderStatusCondition string

const (
	SuccessOIDCProviderStatusCondition                         = OIDCProviderStatusCondition("Success")
	DuplicateOIDCProviderStatusCondition                       = OIDCProviderStatusCondition("Duplicate")
	SameIssuerHostMustUseSameSecretOIDCProviderStatusCondition = OIDCProviderStatusCondition("SameIssuerHostMustUseSameSecret")
	InvalidOIDCProviderStatusCondition                         = OIDCProviderStatusCondition("Invalid")
)

// OIDCProviderTLSSpec is a struct that describes the TLS configuration for an OIDC Provider.
type OIDCProviderTLSSpec struct {
	// SecretName is an optional name of a Secret in the same namespace, of type `kubernetes.io/tls`, which contains
	// the TLS serving certificate for the HTTPS endpoints served by this OIDCProvider. When provided, the TLS Secret
	// named here must contain keys named `tls.crt` and `tls.key` that contain the certificate and private key to use
	// for TLS.
	//
	// Server Name Indication (SNI) is an extension to the Transport Layer Security (TLS) supported by all major browsers.
	//
	// SecretName is required if you would like to use different TLS certificates for issuers of different hostnames.
	// SNI requests do not include port numbers, so all issuers with the same DNS hostname must use the same
	// SecretName value even if they have different port numbers.
	//
	// SecretName is not required when you would like to use only the HTTP endpoints (e.g. when terminating TLS at an
	// Ingress). It is also not required when you would like all requests to this OIDC Provider's HTTPS endpoints to
	// use the default TLS certificate, which is configured elsewhere.
	//
	// When your Issuer URL's host is an IP address, then this field is ignored. SNI does not work for IP addresses.
	//
	// +optional
	SecretName string `json:"secretName,omitempty"`
}

// OIDCProviderSpec is a struct that describes an OIDC Provider.
type OIDCProviderSpec struct {
	// Issuer is the OIDC Provider's issuer, per the OIDC Discovery Metadata document, as well as the
	// identifier that it will use for the iss claim in issued JWTs. This field will also be used as
	// the base URL for any endpoints used by the OIDC Provider (e.g., if your issuer is
	// https://example.com/foo, then your authorization endpoint will look like
	// https://example.com/foo/some/path/to/auth/endpoint).
	//
	// See
	// https://openid.net/specs/openid-connect-discovery-1_0.html#rfc.section.3 for more information.
	// +kubebuilder:validation:MinLength=1
	Issuer string `json:"issuer"`

	// TLS configures how this OIDCProvider is served over Transport Layer Security (TLS).
	// +optional
	TLS *OIDCProviderTLSSpec `json:"tls,omitempty"`
}

// OIDCProviderStatus is a struct that describes the actual state of an OIDC Provider.
type OIDCProviderStatus struct {
	// Status holds an enum that describes the state of this OIDC Provider. Note that this Status can
	// represent success or failure.
	// +optional
	Status OIDCProviderStatusCondition `json:"status,omitempty"`

	// Message provides human-readable details about the Status.
	// +optional
	Message string `json:"message,omitempty"`

	// LastUpdateTime holds the time at which the Status was last updated. It is a pointer to get
	// around some undesirable behavior with respect to the empty metav1.Time value (see
	// https://github.com/kubernetes/kubernetes/issues/86811).
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`

	// JWKSSecret holds the name of the secret in which this OIDC Provider's signing/verification keys
	// are stored. If it is empty, then the signing/verification keys are either unknown or they don't
	// exist.
	// +optional
	JWKSSecret corev1.LocalObjectReference `json:"jwksSecret,omitempty"`
}

// OIDCProvider describes the configuration of an OIDC provider.
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OIDCProvider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec of the OIDC provider.
	Spec OIDCProviderSpec `json:"spec"`

	// Status of the OIDC provider.
	Status OIDCProviderStatus `json:"status,omitempty"`
}

// List of OIDCProvider objects.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type OIDCProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []OIDCProvider `json:"items"`
}
