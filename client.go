package kctest

import "sigs.k8s.io/controller-runtime/pkg/client"

func (c *Cluster) Create(obj client.Object, opts ...client.CreateOption) error {
	return c.client.Create(c.ctx, obj, opts...)
}
