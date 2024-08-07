package kctest

import (
	"context"
	"testing"
)

func TestNewEnvTestCluster(t *testing.T) {
	scenarios := []struct {
		name           string
		clusterName    string
		config         Config
		cleanupCluster bool
	}{
		{
			name:           "Function is called. Result: testEnv Cluster is created",
			clusterName:    "test-cluster",
			cleanupCluster: true,
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewCluster(context.Background(), tc.clusterName, tc.config)
			defer func() {
				if err := result.Stop(); err != nil {
					t.Fatalf("Test failed: %s\nFailed to Cleanup Cluster.", tc.name)
				}
			}() // fail safe
			if !tc.cleanupCluster {
				if err == nil {
					t.Fatalf("Test Failed: %s", tc.name)
				}
				return
			}
			res := result.ClientSet().RESTClient().Get().AbsPath("/healthz").Do(context.Background())
			if res.Error() != nil {
				// not possible to reach the server
				t.Fatalf("Test: %s\nCluster could not be reached: %s", tc.name, result.Name())
			}
			if err := result.Stop(); err != nil {
				t.Fatalf("Test failed: %s\nFailed to Cleanup Cluster.", tc.name)
			}
		})
	}
}

func TestNewKindCluster(t *testing.T) {
	scenarios := []struct {
		name           string
		clusterName    string
		config         Config
		clusterCreated bool
	}{
		{
			name:           "An invalid clusterName is provided. Result: Kind Cluster is not created",
			clusterName:    "new_kind_cluster1",
			config:         Config{kindCluster: true},
			clusterCreated: false,
		},
		{
			name:           "A valid clusterName is provided. Result: Kind Cluster is created",
			clusterName:    "new-kind-cluster2",
			config:         Config{kindCluster: true},
			clusterCreated: true,
		},
		{
			name:           "An empty clusterName is provided. Result: Kind Cluster is not created",
			clusterName:    "",
			config:         Config{kindCluster: true},
			clusterCreated: false,
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewCluster(context.Background(), tc.clusterName, tc.config)
			if !tc.clusterCreated {
				if err == nil {
					t.Fatalf("Test Failed: %s", tc.name)
				}
				return
			}
			defer func() { _ = result.Stop() }()
			res := result.ClientSet().RESTClient().Get().AbsPath("/healthz").Do(context.Background())
			if res.Error() != nil {
				// not possible to reach the server
				t.Fatalf("cluster could not be reached: %s", result.Name())
			}
		})
	}
}
