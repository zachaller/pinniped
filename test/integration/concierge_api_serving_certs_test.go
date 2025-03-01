// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	loginv1alpha1 "go.pinniped.dev/generated/latest/apis/concierge/login/v1alpha1"
	"go.pinniped.dev/internal/testutil"
	"go.pinniped.dev/test/library"
)

func TestAPIServingCertificateAutoCreationAndRotation(t *testing.T) {
	env := library.IntegrationEnv(t)
	defaultServingCertResourceName := env.ConciergeAppName + "-api-tls-serving-certificate"

	tests := []struct {
		name          string
		forceRotation func(context.Context, kubernetes.Interface, string) error
	}{
		{
			name: "manual",
			forceRotation: func(
				ctx context.Context,
				kubeClient kubernetes.Interface,
				namespace string,
			) error {
				// Delete the Secret, simulating an end user doing `kubectl delete` to manually ask for an immediate rotation.
				return kubeClient.
					CoreV1().
					Secrets(namespace).
					Delete(ctx, defaultServingCertResourceName, metav1.DeleteOptions{})
			},
		},
		{
			name: "automatic",
			forceRotation: func(
				ctx context.Context,
				kubeClient kubernetes.Interface,
				namespace string,
			) error {
				// Create a cert that is expired - this should force the rotation controller
				// to delete the cert, and therefore the cert should get rotated.
				secret, err := kubeClient.
					CoreV1().
					Secrets(namespace).
					Get(ctx, defaultServingCertResourceName, metav1.GetOptions{})
				if err != nil {
					return err
				}

				secret.Data["tlsCertificateChain"], _, err = createExpiredCertificate()
				if err != nil {
					return err
				}

				_, err = kubeClient.
					CoreV1().
					Secrets(namespace).
					Update(ctx, secret, metav1.UpdateOptions{})
				return err
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			kubeClient := library.NewKubernetesClientset(t)
			aggregatedClient := library.NewAggregatedClientset(t)
			conciergeClient := library.NewConciergeClientset(t)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			apiServiceName := "v1alpha1.login.concierge." + env.APIGroupSuffix

			// Create a testWebhook so we have a legitimate authenticator to pass to the
			// TokenCredentialRequest API.
			testWebhook := library.CreateTestWebhookAuthenticator(ctx, t)

			// Get the initial auto-generated version of the Secret.
			secret, err := kubeClient.CoreV1().Secrets(env.ConciergeNamespace).Get(ctx, defaultServingCertResourceName, metav1.GetOptions{})
			require.NoError(t, err)
			initialCACert := secret.Data["caCertificate"]
			initialPrivateKey := secret.Data["tlsPrivateKey"]
			initialCertChain := secret.Data["tlsCertificateChain"]
			require.NotEmpty(t, initialCACert)
			require.NotEmpty(t, initialPrivateKey)
			require.NotEmpty(t, initialCertChain)
			for k, v := range env.ConciergeCustomLabels {
				require.Equalf(t, v, secret.Labels[k], "expected secret to have label %s: %s", k, v)
			}
			require.Equal(t, env.ConciergeAppName, secret.Labels["app"])

			// Check that the APIService has the same CA.
			apiService, err := aggregatedClient.ApiregistrationV1().APIServices().Get(ctx, apiServiceName, metav1.GetOptions{})
			require.NoError(t, err)
			require.Equal(t, initialCACert, apiService.Spec.CABundle)

			// Force rotation to happen.
			require.NoError(t, test.forceRotation(ctx, kubeClient, env.ConciergeNamespace))

			// Expect that the Secret comes back right away with newly minted certs.
			secretIsRegenerated := func() bool {
				secret, err = kubeClient.CoreV1().Secrets(env.ConciergeNamespace).Get(ctx, defaultServingCertResourceName, metav1.GetOptions{})
				return err == nil
			}
			assert.Eventually(t, secretIsRegenerated, 10*time.Second, 250*time.Millisecond)
			require.NoError(t, err) // prints out the error and stops the test in case of failure
			regeneratedCACert := secret.Data["caCertificate"]
			regeneratedPrivateKey := secret.Data["tlsPrivateKey"]
			regeneratedCertChain := secret.Data["tlsCertificateChain"]
			require.NotEmpty(t, regeneratedCACert)
			require.NotEmpty(t, regeneratedPrivateKey)
			require.NotEmpty(t, regeneratedCertChain)
			require.NotEqual(t, initialCACert, regeneratedCACert)
			require.NotEqual(t, initialPrivateKey, regeneratedPrivateKey)
			require.NotEqual(t, initialCertChain, regeneratedCertChain)
			for k, v := range env.ConciergeCustomLabels {
				require.Equalf(t, v, secret.Labels[k], "expected secret to have label `%s: %s`", k, v)
			}
			require.Equal(t, env.ConciergeAppName, secret.Labels["app"])

			// Expect that the APIService was also updated with the new CA.
			require.Eventually(t, func() bool {
				apiService, err := aggregatedClient.ApiregistrationV1().APIServices().Get(ctx, apiServiceName, metav1.GetOptions{})
				if err != nil {
					t.Logf("get for APIService %q returned error %v", apiServiceName, err)
					return false
				}
				if !bytes.Equal(regeneratedCACert, apiService.Spec.CABundle) {
					t.Logf("CA bundle in APIService %q does not yet have the expected value", apiServiceName)
					return false
				}
				t.Logf("found that APIService %q was updated to expected CA certificate", apiServiceName)
				return true
			}, 10*time.Second, 250*time.Millisecond, "never saw CA certificate rotate to expected value")

			// Check that we can still make requests to the aggregated API through the kube API server,
			// because the kube API server uses these certs when proxying requests to the aggregated API server,
			// so this is effectively checking that the aggregated API server is using these new certs.
			// We ensure that 10 straight requests succeed so that we filter out false positives where a single
			// pod has rotated their cert, but not the other ones sitting behind the service.
			aggregatedAPIWorking := func() bool {
				for i := 0; i < 10; i++ {
					_, err = conciergeClient.LoginV1alpha1().TokenCredentialRequests().Create(ctx, &loginv1alpha1.TokenCredentialRequest{
						TypeMeta:   metav1.TypeMeta{},
						ObjectMeta: metav1.ObjectMeta{},
						Spec:       loginv1alpha1.TokenCredentialRequestSpec{Token: "not a good token", Authenticator: testWebhook},
					}, metav1.CreateOptions{})
					if err != nil {
						break
					}
				}
				// Should have got a success response with an error message inside it complaining about the token value.
				return err == nil
			}

			// our code changes all the certs immediately thus this should be healthy fairly quickly
			// if this starts flaking, check for bugs in our dynamiccertificates.Notifier implementation
			assert.Eventually(t, aggregatedAPIWorking, 30*time.Second, 250*time.Millisecond)
			require.NoError(t, err, "dynamiccertificates.Notifier broken?") // prints out the error and stops the test in case of failure
		})
	}
}

func createExpiredCertificate() ([]byte, []byte, error) {
	return testutil.CreateCertificate(
		time.Now().Add(-24*time.Hour), // notBefore
		time.Now().Add(-time.Hour),    // notAfter
	)
}
