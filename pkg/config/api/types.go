/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package api

// Config contains knobs to setup an instance of pinniped.
type Config struct {
	WebhookConfig WebhookConfigSpec `json:"webhook"`
	DiscoveryInfo DiscoveryInfoSpec `json:"discovery"`
	APIConfig     APIConfigSpec     `json:"api"`
}

// WebhookConfig contains configuration knobs specific to pinniped's use
// of a webhook for token validation.
type WebhookConfigSpec struct {
	// URL contains the URL of the webhook that pinniped will use
	// to validate external credentials.
	URL string `json:"url"`

	// CABundle contains PEM-encoded certificate authority certificates used
	// to validate TLS connections to the WebhookURL.
	CABundle []byte `json:"caBundle"`
}

// DiscoveryInfoSpec contains configuration knobs specific to
// pinniped's publishing of discovery information. These values can be
// viewed as overrides, i.e., if these are set, then pinniped will
// publish these values in its discovery document instead of the ones it finds.
type DiscoveryInfoSpec struct {
	// URL contains the URL at which pinniped can be contacted.
	URL *string `json:"url,omitempty"`
}

// APIConfigSpec contains configuration knobs for the pinniped API.
//nolint: golint
type APIConfigSpec struct {
	ServingCertificateConfig ServingCertificateConfigSpec `json:"servingCertificate"`
}

// ServingCertificateConfigSpec contains the configuration knobs for the API's
// serving certificate, i.e., the x509 certificate that it uses for the server
// certificate in inbound TLS connections.
type ServingCertificateConfigSpec struct {
	// DurationSeconds is the validity period, in seconds, of the API serving
	// certificate. By default, the serving certificate is issued for 31536000
	// seconds (1 year).
	DurationSeconds *int64 `json:"durationSeconds,omitempty"`

	// RenewBeforeSeconds is the period of time, in seconds, that pinniped will
	// wait before rotating the serving certificate. This period of time starts
	// upon issuance of the serving certificate. This must be less than
	// DurationSeconds. By default, pinniped begins rotation after 23328000
	// seconds (about 9 months).
	RenewBeforeSeconds *int64 `json:"renewBeforeSeconds,omitempty"`
}
