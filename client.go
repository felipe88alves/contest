package kctest

import "sigs.k8s.io/controller-runtime/pkg/client"

func (c *Cluster) Create(obj client.Object, opts ...client.CreateOption) error {
	return c.client.Create(c.ctx, obj, opts...)
}

func (c *Cluster) Get(obj client.Object, opts ...client.GetOption) func() error {
	key := client.ObjectKeyFromObject(obj)
	return func() error {
		return c.client.Get(c.ctx, key, obj, opts...)
	}
}
