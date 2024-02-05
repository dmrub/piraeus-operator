package podpatcher_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/piraeusdatastore/piraeus-operator/v2/pkg/podpatcher"
)

var (
	ExamplePod1 = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "container-1", Image: "image-1"},
				{Name: "container-2", Image: "image-2"},
			},
			InitContainers: []corev1.Container{
				{Name: "init-1", Image: "image-init-1"},
				{Name: "init-2", Image: "image-init-2"},
			},
		},
	}
	ExamplePod2 = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "container-1", Image: "image-1"},
				{Name: "container-3", Image: "image-3"},
			},
			InitContainers: []corev1.Container{
				{Name: "init-1", Image: "image-init-1"},
				{Name: "init-2", Image: "image-init-2"},
			},
		},
	}
	ExamplePod3 = &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "container-1", Image: "image-1"},
				{Name: "container-2", Image: "image-2"},
			},
			InitContainers: []corev1.Container{
				{Name: "init-1", Image: "image-init-1"},
				{Name: "init-2", Image: "image-init-3"},
			},
		},
	}
)

func TestEqualImages(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		a        *corev1.Pod
		b        *corev1.Pod
		expected bool
	}{
		{
			name:     "same-obj",
			a:        ExamplePod1,
			b:        ExamplePod1,
			expected: true,
		},
		{
			name:     "diff-container",
			a:        ExamplePod1,
			b:        ExamplePod2,
			expected: false,
		},
		{
			name:     "diff-init",
			a:        ExamplePod1,
			b:        ExamplePod3,
			expected: false,
		},
	}

	for i := range testcases {
		tcase := &testcases[i]
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			actual := podpatcher.EqualImages(tcase.a, tcase.b)
			assert.Equal(t, tcase.expected, actual)
		})
	}
}

func TestPodNeedsRestart(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name        string
		pod         *corev1.Pod
		expectation map[corev1.RestartPolicy]bool
	}{
		{
			name: "restart-running",
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodRunning,
				},
			},
			expectation: map[corev1.RestartPolicy]bool{
				corev1.RestartPolicyAlways:    false,
				corev1.RestartPolicyOnFailure: false,
				corev1.RestartPolicyNever:     false,
			},
		},
		{
			name: "restart-failed",
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodFailed,
				},
			},
			expectation: map[corev1.RestartPolicy]bool{
				corev1.RestartPolicyAlways:    true,
				corev1.RestartPolicyOnFailure: true,
				corev1.RestartPolicyNever:     false,
			},
		},
		{
			name: "restart-success",
			pod: &corev1.Pod{
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
				},
			},
			expectation: map[corev1.RestartPolicy]bool{
				corev1.RestartPolicyAlways:    true,
				corev1.RestartPolicyOnFailure: false,
				corev1.RestartPolicyNever:     false,
			},
		},
	}

	for i := range testcases {
		tcase := &testcases[i]
		t.Run(tcase.name, func(t *testing.T) {
			t.Parallel()

			for p, e := range tcase.expectation {
				policy := p
				expected := e
				t.Run(string(policy), func(t *testing.T) {
					t.Parallel()

					actual := podpatcher.PodNeedsRestart(tcase.pod, policy)
					assert.Equal(t, expected, actual)
				})
			}
		})
	}
}

// NB: sadly, the fake client does not support "Apply" patches, so we can't test podpatcher.Patch easily