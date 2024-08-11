package kctest

import (
	"context"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func setupConfigMap(name, ns string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm" + name,
			Namespace: ns,
		}}
}

func setupNamespace(ns string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}
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
	defer t.Cleanup(func() { _ = c.Stop() })

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := new(corev1.ConfigMap)
			if tc.setupResource {
				cm = setupConfigMap(tc.resourceName, "default")
			}

			err = c.Create(cm)
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
		})
	}
}

func TestGet(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Get succeeds.",
			resourceName:  "1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Get fails.",
			resourceName:  "2",
			setupResource: false,
			expectedFail:  true,
		},
	}
	c, err := setupK8sClient()
	if err != nil {
		t.Fatalf("failed to setup kind cluster: %s", err.Error())
	}
	defer t.Cleanup(func() { _ = c.Stop() })
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := new(corev1.ConfigMap)
			if tc.setupResource {
				cm = setupConfigMap(tc.resourceName, "default")
				if err = c.Create(cm); err != nil {
					t.Fatalf("setup failed. Failed to create resource: %s", err.Error())
				}
			}

			err := c.Get(cm)()
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
		})
	}
}

func TestList(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  []string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resources exist. Get List succeeds.",
			resourceName:  []string{"1", "3"},
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resources do not exist. Get List fails.",
			resourceName:  []string{"2", "4"},
			setupResource: false,
			expectedFail:  true,
		},
	}
	c, err := setupK8sClient()
	if err != nil {
		t.Fatalf("failed to setup kind cluster: %s", err.Error())
	}
	defer t.Cleanup(func() { _ = c.Stop() })
	var nsId int
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// setup NS
			ns := setupNamespace(fmt.Sprintf("ns-%d", nsId))
			nsId++
			if err := c.Create(ns); err != nil {
				t.Fatalf("setup failed. Failed to create namespace: %s", err.Error())
			}
			if tc.setupResource {
				// setup CM
				for _, r := range tc.resourceName {
					cm := setupConfigMap(r, ns.GetName())
					if err := c.Create(cm); err != nil {
						t.Fatalf("setup failed. Failed to create resource: %s", err.Error())
					}
				}
			}

			cmList := new(corev1.ConfigMapList)
			opts := []client.ListOption{client.InNamespace(ns.GetName())}
			err := c.List(cmList, opts...)()
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				return
			}
			if !tc.expectedFail {
				if len(cmList.Items) != len(tc.resourceName) {
					t.Fatalf("Length of Resources in List don't match: Expected: %d, Actual: %d", len(tc.resourceName), len(cmList.Items))
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Delete succeeds.",
			resourceName:  "1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Delete fails.",
			resourceName:  "2",
			setupResource: false,
			expectedFail:  true,
		},
	}
	c, err := setupK8sClient()
	if err != nil {
		t.Fatalf("failed to setup kind cluster: %s", err.Error())
	}
	defer t.Cleanup(func() { _ = c.Stop() })
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cm := new(corev1.ConfigMap)
			if tc.setupResource {
				cm = setupConfigMap(tc.resourceName, "default")
				if err = c.Create(cm); err != nil {
					t.Fatalf("setup failed. Failed to create resource: %s", err.Error())
				}
			}

			err := c.Delete(cm)()
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
		})
	}
}
