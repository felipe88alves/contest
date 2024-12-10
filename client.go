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

func (c *Cluster) List(objList client.ObjectList, opts ...client.ListOption) func() error {
	return func() error {
		return c.client.List(c.ctx, objList, opts...)
	}
}

func (c *Cluster) Object(obj client.Object) func() (client.Object, error) {
	key := client.ObjectKeyFromObject(obj)
	return func() (client.Object, error) {
		err := c.client.Get(c.ctx, key, obj)
		return obj, err
	}
}

func (c *Cluster) ObjectList(objList client.ObjectList, opts ...client.ListOption) func() (client.ObjectList, error) {
	return func() (client.ObjectList, error) {
		err := c.client.List(c.ctx, objList, opts...)
		return objList, err
	}
}

func (c *Cluster) Update(obj client.Object, update func() error, opts ...client.UpdateOption) func() error {
	key := client.ObjectKeyFromObject(obj)
	return func() error {
		if err := c.client.Get(c.ctx, key, obj); err != nil {
			return err
		}

		if err := update(); err != nil {
			return err
		}

		return c.client.Update(c.ctx, obj, opts...)
	}
}

func (c *Cluster) UpdateStatus(obj client.Object, update func() error, opts ...client.SubResourceUpdateOption) func() error {
	key := client.ObjectKeyFromObject(obj)
	return func() error {
		if err := c.client.Get(c.ctx, key, obj); err != nil {
			return err
		}

		if err := update(); err != nil {
			return err
		}

		return c.client.Status().Update(c.ctx, obj, opts...)
	}
}
