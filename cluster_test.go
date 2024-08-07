package kctest

import (
	"context"
	"testing"
)

func TestNewCluster(t *testing.T) {
	scenarios := []struct {
		name           string
		clusterName    string
		cleanupCluster bool
	}{
		{
			name:           "An invalid clusterName is provided. Result: testEnv Cluster is not created",
			clusterName:    "test cluster",
			cleanupCluster: false,
		},
	}

	for _, tc := range scenarios {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewCluster(tc.clusterName)
			defer func() { _ = result.Stop() }()
			if tc.cleanupCluster {
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
		})
	}
}
