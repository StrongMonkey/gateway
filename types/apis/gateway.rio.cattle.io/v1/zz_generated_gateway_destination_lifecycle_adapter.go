package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type GatewayDestinationLifecycle interface {
	Create(obj *GatewayDestination) (runtime.Object, error)
	Remove(obj *GatewayDestination) (runtime.Object, error)
	Updated(obj *GatewayDestination) (runtime.Object, error)
}

type gatewayDestinationLifecycleAdapter struct {
	lifecycle GatewayDestinationLifecycle
}

func (w *gatewayDestinationLifecycleAdapter) HasCreate() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasCreate()
}

func (w *gatewayDestinationLifecycleAdapter) HasFinalize() bool {
	o, ok := w.lifecycle.(lifecycle.ObjectLifecycleCondition)
	return !ok || o.HasFinalize()
}

func (w *gatewayDestinationLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*GatewayDestination))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gatewayDestinationLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*GatewayDestination))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *gatewayDestinationLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*GatewayDestination))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewGatewayDestinationLifecycleAdapter(name string, clusterScoped bool, client GatewayDestinationInterface, l GatewayDestinationLifecycle) GatewayDestinationHandlerFunc {
	adapter := &gatewayDestinationLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *GatewayDestination) (runtime.Object, error) {
		newObj, err := syncFn(key, obj)
		if o, ok := newObj.(runtime.Object); ok {
			return o, err
		}
		return nil, err
	}
}
