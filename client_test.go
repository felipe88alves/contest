package kctest

import (
	"context"
	"fmt"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/api/meta"
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
			Name:      name,
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
			resourceName:  "cm-1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource is not setup. Create fails.",
			resourceName:  "cm-2",
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

func TestDelete(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Delete succeeds.",
			resourceName:  "cm-1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Delete fails.",
			resourceName:  "cm-2",
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

func TestGet(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Get succeeds.",
			resourceName:  "cm-1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Get fails.",
			resourceName:  "cm-2",
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
			resourceName:  []string{"cm-1", "cm-3"},
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resources do not exist. Get List fails.",
			resourceName:  []string{"cm-2", "cm-4"},
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

func TestObject(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Object is captured and contains appropriate data.",
			resourceName:  "cm-1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Object is not captured and fails.",
			resourceName:  "cm-2",
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

			obj, err := c.Object(cm)()
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
			if !tc.expectedFail && (obj.GetName() != tc.resourceName || obj.GetNamespace() != "default") {
				t.Fatalf("Object does not contain correct value: Expected Namespace/Name: %s/%s, Actual: %s/%s", "default", tc.resourceName, obj.GetNamespace(), obj.GetName())
			}
		})
	}
}

func TestObjectList(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  []string
		setupResource bool
	}{
		{
			name:          "Resources exist. Object List is captured and contain appropriate data.",
			resourceName:  []string{"cm-1", "cm-3"},
			setupResource: true,
		},
		{
			name:          "Resources do not exist. Objects are not captured and fails.",
			resourceName:  []string{"cm-2", "cm-4"},
			setupResource: false,
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
			objList, err := c.ObjectList(cmList, opts...)()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			objs, err := meta.ExtractList(objList)
			if err != nil {
				t.Fatal("Failed to access Objects in List. Failed to extract items from ObjectList")
			}
			for _, obj := range objs {
				o, ok := obj.(client.Object)
				if !ok {
					t.Fatal("Failed to access Objects in List. Type assertion failed")
				}
				if !slices.Contains(tc.resourceName, o.GetName()) || o.GetNamespace() != ns.GetName() {
					t.Fatalf("Object in ObjectList does not contain correct value: Expected Namespace/Name: %s/%s, Actual: %s/%s", ns.GetName(), tc.resourceName, o.GetNamespace(), o.GetName())
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	scenarios := []struct {
		name          string
		resourceName  string
		setupResource bool
		expectedFail  bool
	}{
		{
			name:          "Resource exists. Object is updated and contains appropriate data.",
			resourceName:  "cm-1",
			setupResource: true,
			expectedFail:  false,
		},
		{
			name:          "Resource does not exist. Object is not found and is not updated.",
			resourceName:  "cm-2",
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

			err := c.Update(cm, func() error {
				cm.SetLabels(map[string]string{
					"state": "updated",
				})
				return nil
			})()
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Expected to fail: %t, Result: %t, Error: %s", tc.expectedFail, errResult, err.Error())
			}
			cmLabels := cm.GetLabels()
			stateLabel, ok := cmLabels["state"]
			if !tc.expectedFail && !ok && stateLabel != "updated" {
				t.Fatalf("Object does not contain correct value: Expected Updated Label: \"state\", %q, Actual: %s", stateLabel, cmLabels)
			}
		})
	}
}
