// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package authenticator contains helper code for dealing with *Authenticator CRDs.
package authenticator

import (
	"encoding/base64"

	auth1alpha1 "go.pinniped.dev/generated/latest/apis/concierge/authentication/v1alpha1"
)

// Closer is a type that can be closed idempotently.
//
// This type is slightly different from io.Closer, because io.Closer can return an error and is not
// necessarily idempotent.
type Closer interface {
	Close()
}

// CABundle returns a PEM-encoded CA bundle from the provided spec. If the provided spec is nil, a
// nil CA bundle will be returned. If the provided spec contains a CA bundle that is not properly
// encoded, an error will be returned.
func CABundle(spec *auth1alpha1.TLSSpec) ([]byte, error) {
	if spec == nil {
		return nil, nil
	}
	return base64.StdEncoding.DecodeString(spec.CertificateAuthorityData)
}
