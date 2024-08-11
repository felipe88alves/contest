package kctest

import "sigs.k8s.io/controller-runtime/pkg/client"

func (c *Cluster) Create(obj client.Object, opts ...client.CreateOption) error {
	return c.client.Create(c.ctx, obj, opts...)
}

func (c *Cluster) Delete(obj client.Object, opts ...client.DeleteOption) func() error {
	return func() error {
		if err := c.Get(obj)(); err != nil {
			return err
		}
		return c.client.Delete(c.ctx, obj, opts...)
	}
}

func (c *Cluster) Get(obj client.Object, opts ...client.GetOption) func() error {
	key := client.ObjectKeyFromObject(obj)
	return func() error {
		return c.client.Get(c.ctx, key, obj, opts...)
	}
}

func (c *Cluster) List(obj client.ObjectList, opts ...client.ListOption) func() error {
	return func() error {
		return c.client.List(c.ctx, obj, opts...)
	}
}
