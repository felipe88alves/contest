package kctest

import (
	"context"
	"testing"
)

func TestKindClusterCreate(t *testing.T) {
	scenarios := []struct {
		name        string
		clusterName string
		config      Config
		shouldFail  bool
	}{
		{
			name:        "Cluster name is valid. Result: Kind cluster is created.",
			clusterName: "test-cluster1",
			config:      Config{kindCluster: true},
			shouldFail:  false,
		},
		{
			name:        "Cluster name is invalid. Result: Kind cluster is not created.",
			clusterName: "test_cluster",
			config:      Config{kindCluster: true},
			shouldFail:  true,
		},
		{
			name:        "Cluster config set kind to false. Result: Kind cluster is not created.",
			clusterName: "test-cluster2",
			config:      Config{kindCluster: false},
			shouldFail:  false,
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := kindClusterCreate(context.Background(), tc.clusterName, tc.config.kindCluster)
			defer func() { _ = kindClusterDelete(tc.clusterName) }()

			errResult := (err != nil)
			if errResult != tc.shouldFail {
				t.Fatalf("Test Failed: %s\nExpected to fail: %t, Result: %t, Error: %s", tc.name, tc.shouldFail, errResult, err)
			}
		})
	}
}

func TestKindClusterDelete(t *testing.T) {
	scenarios := []struct {
		name         string
		clusterName  string
		setupCluster bool
		expectedFail bool
	}{
		{
			name:         "clusterName not found. Result: Delete fails.",
			clusterName:  "kind-delete-cluster1",
			setupCluster: false,
			expectedFail: true,
		},
		{
			name:         "clusterName found. Result: Cluster deleted .",
			clusterName:  "kind-delete-cluster2",
			setupCluster: true,
			expectedFail: false,
		},
	}
	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.setupCluster {
				if err := kindClusterCreate(context.Background(), tc.clusterName, true); err != nil {
					t.Fatal("failed to setup Kind Cluster")
				}
			}

			err := kindClusterDelete(tc.clusterName)
			errResult := (err != nil)
			if errResult != tc.expectedFail {
				t.Fatalf("Test Failed: %s\nExpected to fail: %t, Result: %t, Error: %s", tc.name, tc.expectedFail, errResult, err)
			}

		})
	}
}
