// Copyright 2020-2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/apimachinery/pkg/util/wait"

	conciergev1alpha "go.pinniped.dev/generated/latest/apis/concierge/config/v1alpha1"
	"go.pinniped.dev/test/library"
)

const (
	kubeCertAgentLabelSelector = "kube-cert-agent.pinniped.dev=true"
)

func TestKubeCertAgent(t *testing.T) {
	env := library.IntegrationEnv(t).WithCapability(library.ClusterSigningKeyIsAvailable)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeClient := library.NewKubernetesClientset(t)

	// Get the current number of kube-cert-agent pods.
	//
	// We can pretty safely assert there should be more than 1, since there should be a
	// kube-cert-agent pod per kube-controller-manager pod, and there should probably be at least
	// 1 kube-controller-manager for this to be a working kube API.
	originalAgentPods, err := kubeClient.CoreV1().Pods(env.ConciergeNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: kubeCertAgentLabelSelector,
	})
	require.NoError(t, err)
	require.NotEmpty(t, originalAgentPods.Items)
	sortPods(originalAgentPods)

	for _, agentPod := range originalAgentPods.Items {
		// All agent pods should contain all custom labels
		for k, v := range env.ConciergeCustomLabels {
			require.Equalf(t, v, agentPod.Labels[k], "expected agent pod to have label `%s: %s`", k, v)
		}
		require.Equal(t, env.ConciergeAppName, agentPod.Labels["app"])
	}

	agentPodsReconciled := func() bool {
		var currentAgentPods *corev1.PodList
		currentAgentPods, err = kubeClient.CoreV1().Pods(env.ConciergeNamespace).List(ctx, metav1.ListOptions{
			LabelSelector: kubeCertAgentLabelSelector,
		})

		if err != nil {
			return false
		}

		if len(originalAgentPods.Items) != len(currentAgentPods.Items) {
			err = fmt.Errorf(
				"original agent pod len != current agent pod len: %s",
				diff.ObjectDiff(originalAgentPods.Items, currentAgentPods.Items),
			)
			return false
		}

		sortPods(currentAgentPods)
		for i := range originalAgentPods.Items {
			if !equality.Semantic.DeepEqual(
				originalAgentPods.Items[i].Spec,
				currentAgentPods.Items[i].Spec,
			) {
				err = fmt.Errorf(
					"original agent pod != current agent pod: %s",
					diff.ObjectDiff(originalAgentPods.Items[i].Spec, currentAgentPods.Items[i].Spec),
				)
				return false
			}
		}

		return true
	}

	t.Run("reconcile on update", func(t *testing.T) {
		// Ensure that the next test will start from a known state.
		defer ensureKubeCertAgentSteadyState(t, agentPodsReconciled)

		// Update the image of the first pod. The controller should see it, and flip it back.
		//
		// Note that we update the toleration field here because it is the only field, currently, that
		// 1) we are allowed to update on a running pod AND 2) the kube-cert-agent controllers care
		// about.
		updatedAgentPod := originalAgentPods.Items[0].DeepCopy()
		updatedAgentPod.Spec.Tolerations = append(
			updatedAgentPod.Spec.Tolerations,
			corev1.Toleration{Key: "fake-toleration"},
		)
		_, err = kubeClient.CoreV1().Pods(env.ConciergeNamespace).Update(ctx, updatedAgentPod, metav1.UpdateOptions{})
		require.NoError(t, err)

		// Make sure the original pods come back.
		assert.Eventually(t, agentPodsReconciled, 10*time.Second, 250*time.Millisecond)
		require.NoError(t, err)
	})

	t.Run("reconcile on delete", func(t *testing.T) {
		// Ensure that the next test will start from a known state.
		defer ensureKubeCertAgentSteadyState(t, agentPodsReconciled)

		// Delete the first pod. The controller should see it, and flip it back.
		err = kubeClient.
			CoreV1().
			Pods(env.ConciergeNamespace).
			Delete(ctx, originalAgentPods.Items[0].Name, metav1.DeleteOptions{})
		require.NoError(t, err)

		// Make sure the original pods come back.
		assert.Eventually(t, agentPodsReconciled, 10*time.Second, 250*time.Millisecond)
		require.NoError(t, err)
	})

	// Because the above tests have purposefully put the kube cert issuer strategy into a broken
	// state, wait for it to become healthy again before moving on to other integration tests,
	// otherwise those tests would be polluted by this test and would have to wait for the
	// strategy to become successful again.
	library.RequireEventuallyWithoutError(t, func() (bool, error) {
		adminConciergeClient := library.NewConciergeClientset(t)
		credentialIssuer, err := adminConciergeClient.ConfigV1alpha1().CredentialIssuers().Get(ctx, credentialIssuerName(env), metav1.GetOptions{})
		if err != nil || credentialIssuer.Status.Strategies == nil {
			t.Log("Did not find any CredentialIssuer with any strategies")
			return false, nil // didn't find it, but keep trying
		}
		for _, strategy := range credentialIssuer.Status.Strategies {
			// There will be other strategy types in the list, so ignore those.
			if strategy.Type == conciergev1alpha.KubeClusterSigningCertificateStrategyType && strategy.Status == conciergev1alpha.SuccessStrategyStatus { //nolint:nestif
				if strategy.Frontend == nil {
					return false, fmt.Errorf("did not find a Frontend") // unexpected, fail the test
				}
				return true, nil // found it, continue the test!
			}
		}
		t.Log("Did not find any successful KubeClusterSigningCertificate strategy on CredentialIssuer")
		return false, nil // didn't find it, but keep trying
	}, 3*time.Minute, 3*time.Second)
}

func ensureKubeCertAgentSteadyState(t *testing.T, agentPodsReconciled func() bool) {
	t.Helper()

	const wantSteadyStateSnapshots = 3
	var steadyStateSnapshots int
	require.NoError(t, wait.Poll(250*time.Millisecond, 30*time.Second, func() (bool, error) {
		if agentPodsReconciled() {
			steadyStateSnapshots++
		} else {
			steadyStateSnapshots = 0
		}
		return steadyStateSnapshots == wantSteadyStateSnapshots, nil
	}))
}

func sortPods(pods *corev1.PodList) {
	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].Name < pods.Items[j].Name
	})
}
