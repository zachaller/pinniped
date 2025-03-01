// Copyright 2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package dynamiccert

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apiserver/pkg/server/dynamiccertificates"
	"k8s.io/apiserver/pkg/storage/names"

	"go.pinniped.dev/internal/certauthority"
	"go.pinniped.dev/test/library"
)

func TestProviderWithDynamicServingCertificateController(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		f    func(t *testing.T, ca Provider, certKey Private) (wantClientCASubjects [][]byte, wantCerts []tls.Certificate)
	}{
		{
			name: "no-op leave everything alone",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				pool := x509.NewCertPool()
				ok := pool.AppendCertsFromPEM(ca.CurrentCABundleContent())
				require.True(t, ok, "should have valid non-empty CA bundle")

				certPEM, keyPEM := certKey.CurrentCertKeyContent()
				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				return pool.Subjects(), []tls.Certificate{cert}
			},
		},
		{
			name: "unset the CA",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				ca.UnsetCertKeyContent()

				certPEM, keyPEM := certKey.CurrentCertKeyContent()
				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				return nil, []tls.Certificate{cert}
			},
		},
		{
			name: "unset the serving cert - still serves the old content",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				pool := x509.NewCertPool()
				ok := pool.AppendCertsFromPEM(ca.CurrentCABundleContent())
				require.True(t, ok, "should have valid non-empty CA bundle")

				certPEM, keyPEM := certKey.CurrentCertKeyContent()
				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				certKey.UnsetCertKeyContent()

				return pool.Subjects(), []tls.Certificate{cert}
			},
		},
		{
			name: "change to a new CA",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				// use unique names for all CAs to make sure the pool subjects are different
				newCA, err := certauthority.New(names.SimpleNameGenerator.GenerateName("new-ca"), time.Hour)
				require.NoError(t, err)
				caKey, err := newCA.PrivateKeyToPEM()
				require.NoError(t, err)
				err = ca.SetCertKeyContent(newCA.Bundle(), caKey)
				require.NoError(t, err)

				certPEM, keyPEM := certKey.CurrentCertKeyContent()
				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				return newCA.Pool().Subjects(), []tls.Certificate{cert}
			},
		},
		{
			name: "change to new serving cert",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				// use unique names for all CAs to make sure the pool subjects are different
				newCA, err := certauthority.New(names.SimpleNameGenerator.GenerateName("new-ca"), time.Hour)
				require.NoError(t, err)

				certPEM, keyPEM, err := newCA.IssueServerCertPEM(nil, []net.IP{net.ParseIP("127.0.0.2")}, time.Hour)
				require.NoError(t, err)

				err = certKey.SetCertKeyContent(certPEM, keyPEM)
				require.NoError(t, err)

				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				pool := x509.NewCertPool()
				ok := pool.AppendCertsFromPEM(ca.CurrentCABundleContent())
				require.True(t, ok, "should have valid non-empty CA bundle")

				return pool.Subjects(), []tls.Certificate{cert}
			},
		},
		{
			name: "change both CA and serving cert",
			f: func(t *testing.T, ca Provider, certKey Private) ([][]byte, []tls.Certificate) {
				// use unique names for all CAs to make sure the pool subjects are different
				newCA, err := certauthority.New(names.SimpleNameGenerator.GenerateName("new-ca"), time.Hour)
				require.NoError(t, err)

				certPEM, keyPEM, err := newCA.IssueServerCertPEM(nil, []net.IP{net.ParseIP("127.0.0.3")}, time.Hour)
				require.NoError(t, err)

				err = certKey.SetCertKeyContent(certPEM, keyPEM)
				require.NoError(t, err)

				cert, err := tls.X509KeyPair(certPEM, keyPEM)
				require.NoError(t, err)

				// use unique names for all CAs to make sure the pool subjects are different
				newOtherCA, err := certauthority.New(names.SimpleNameGenerator.GenerateName("new-other-ca"), time.Hour)
				require.NoError(t, err)
				caKey, err := newOtherCA.PrivateKeyToPEM()
				require.NoError(t, err)
				err = ca.SetCertKeyContent(newOtherCA.Bundle(), caKey)
				require.NoError(t, err)

				return newOtherCA.Pool().Subjects(), []tls.Certificate{cert}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// use unique names for all CAs to make sure the pool subjects are different
			ca, err := certauthority.New(names.SimpleNameGenerator.GenerateName("ca"), time.Hour)
			require.NoError(t, err)
			caKey, err := ca.PrivateKeyToPEM()
			require.NoError(t, err)
			caContent := NewCA("ca")
			err = caContent.SetCertKeyContent(ca.Bundle(), caKey)
			require.NoError(t, err)

			cert, key, err := ca.IssueServerCertPEM(nil, []net.IP{net.ParseIP("127.0.0.1")}, time.Hour)
			require.NoError(t, err)
			certKeyContent := NewServingCert("cert-key")
			err = certKeyContent.SetCertKeyContent(cert, key)
			require.NoError(t, err)

			tlsConfig := &tls.Config{
				MinVersion: tls.VersionTLS12,
				NextProtos: []string{"h2", "http/1.1"},
				ClientAuth: tls.RequestClientCert,
			}

			dynamicCertificateController := dynamiccertificates.NewDynamicServingCertificateController(
				tlsConfig,
				caContent,
				certKeyContent,
				nil, // we do not care about SNI
				nil, // we do not care about events
			)

			caContent.AddListener(dynamicCertificateController)
			certKeyContent.AddListener(dynamicCertificateController)

			err = dynamicCertificateController.RunOnce()
			require.NoError(t, err)

			stopCh := make(chan struct{})
			defer close(stopCh)
			go dynamicCertificateController.Run(1, stopCh)

			tlsConfig.GetConfigForClient = dynamicCertificateController.GetConfigForClient

			wantClientCASubjects, wantCerts := tt.f(t, caContent, certKeyContent)

			var lastTLSConfig *tls.Config

			// it will take some time for the controller to catch up
			err = wait.PollImmediate(time.Second, 30*time.Second, func() (bool, error) {
				actualTLSConfig, err := tlsConfig.GetConfigForClient(&tls.ClientHelloInfo{ServerName: "force-standard-sni"})
				if err != nil {
					return false, err
				}

				lastTLSConfig = actualTLSConfig

				return reflect.DeepEqual(wantClientCASubjects, poolSubjects(actualTLSConfig.ClientCAs)) &&
					reflect.DeepEqual(wantCerts, actualTLSConfig.Certificates), nil
			})

			if err != nil && lastTLSConfig != nil {
				// for debugging failures
				t.Log("diff between client CAs:\n", cmp.Diff(
					library.Sdump(wantClientCASubjects),
					library.Sdump(poolSubjects(lastTLSConfig.ClientCAs)),
				))
				t.Log("diff between serving certs:\n", cmp.Diff(
					library.Sdump(wantCerts),
					library.Sdump(lastTLSConfig.Certificates),
				))
			}
			require.NoError(t, err)
		})
	}
}

func poolSubjects(pool *x509.CertPool) [][]byte {
	if pool == nil {
		return nil
	}
	return pool.Subjects()
}
