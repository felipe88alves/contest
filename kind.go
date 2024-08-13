package kctest

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func kindClusterCreate(ctx context.Context, name string, isKindCluster bool) error {
	if !isKindCluster {
		return nil
	}

	if err := dockerCheck(ctx); err != nil {
		return fmt.Errorf("verify if docker is installed. Error: %w", err)
	}
	cmd := exec.CommandContext(ctx, "kind", "create", "cluster", "--name", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create Kind cluster: %s. Verify if kind is installed", output)
	}
	return nil
}

func kindClusterDelete(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "kind", "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get Kind clusters: %s. Verify if kind is installed", output)
	}
	clusterSlice := strings.Split(string(output), "\n")
	if slices.Contains(clusterSlice[:len(clusterSlice)-1], name) {
		cmd := exec.CommandContext(ctx, "kind", "delete", "cluster", "--name", name)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete Kind cluster: %s", output)
		}
		return nil
	}
	return fmt.Errorf("Cluster %s not found", name)
}

func dockerCheck(ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	_, err = cli.ContainerList(ctx, container.ListOptions{All: true})
	return err
}
