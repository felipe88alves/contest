package kctest

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	ginkgo "github.com/onsi/ginkgo/v2"
)

type Cluster struct {
	name        string
	env         *envtest.Environment
	client      client.Client
	clientSet   *kubernetes.Clientset
	kindCluster bool
	ctx         context.Context
}

type Config struct {
	KindCluster bool
}

func NewCluster(ctx context.Context, name string, config Config) (*Cluster, error) {
	logf.SetLogger(zap.New(zap.WriteTo(ginkgo.GinkgoWriter), zap.UseDevMode(true)))
	if name == "" {
		return nil, fmt.Errorf("a cluster name must be provided")
	}
	if err := kindClusterCreate(ctx, name, config.KindCluster); err != nil {
		return nil, err
	}

	testEnv := &envtest.Environment{
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("bin", "k8s",
			fmt.Sprintf("1.30.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	cfg, err := testEnv.Start()
	if err != nil {
		return nil, err
	}

	cli, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	cluster := &Cluster{
		name:        name,
		env:         testEnv,
		client:      cli,
		clientSet:   clientSet,
		ctx:         ctx,
		kindCluster: config.KindCluster,
	}
	return cluster, nil

}

func (c *Cluster) Stop() error {
	if err := c.env.Stop(); err != nil {
		return err
	}
	if c.kindCluster {
		return kindClusterDelete(c.name)
	}
	return nil
}

func (c *Cluster) Client() client.Client {
	return c.client
}

func (c *Cluster) ClientSet() *kubernetes.Clientset {
	return c.clientSet
}

func (c *Cluster) Name() string {
	return c.name
}
