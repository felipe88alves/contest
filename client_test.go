package kctest

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func setupK8sClient() (*Cluster, error) {
	if err := corev1.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	c, err := NewCluster(context.Background(), "test", Config{})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func setupConfigMap(name string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm" + name,
			Namespace: "default",
		}}
}

func TestCreate(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource is setup. Create succeeds.",
			resourceName:  "1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource is not setup. Create fails.",
			resourceName:  "2",
			setupResource: false,
			expectedFail:  true,
		},
	}
	c, err := setupK8sClient()
	if err != nil {
		t.Fatalf("failed to setup kind cluster: %s", err.Error())
	}
	defer func() { _ = c.Stop() }()
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			cm := new(corev1.ConfigMap)
			if tc.setupResource {
				cm = setupConfigMap(tc.resourceName)
			}

			err = c.Create(cm)
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
		})
	}
}
