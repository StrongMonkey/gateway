package v1

import (
	"context"

	"github.com/rancher/norman/controller"
	"github.com/rancher/norman/objectclient"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	GatewayDestinationGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "GatewayDestination",
	}
	GatewayDestinationResource = metav1.APIResource{
		Name:         "gatewaydestinations",
		SingularName: "gatewaydestination",
		Namespaced:   true,

		Kind: GatewayDestinationGroupVersionKind.Kind,
	}
)

type GatewayDestinationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GatewayDestination
}

type GatewayDestinationHandlerFunc func(key string, obj *GatewayDestination) (runtime.Object, error)

type GatewayDestinationChangeHandlerFunc func(obj *GatewayDestination) (runtime.Object, error)

type GatewayDestinationLister interface {
	List(namespace string, selector labels.Selector) (ret []*GatewayDestination, err error)
	Get(namespace, name string) (*GatewayDestination, error)
}

type GatewayDestinationController interface {
	Generic() controller.GenericController
	Informer() cache.SharedIndexInformer
	Lister() GatewayDestinationLister
	AddHandler(ctx context.Context, name string, handler GatewayDestinationHandlerFunc)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, handler GatewayDestinationHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type GatewayDestinationInterface interface {
	ObjectClient() *objectclient.ObjectClient
	Create(*GatewayDestination) (*GatewayDestination, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GatewayDestination, error)
	Get(name string, opts metav1.GetOptions) (*GatewayDestination, error)
	Update(*GatewayDestination) (*GatewayDestination, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*GatewayDestinationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() GatewayDestinationController
	AddHandler(ctx context.Context, name string, sync GatewayDestinationHandlerFunc)
	AddLifecycle(ctx context.Context, name string, lifecycle GatewayDestinationLifecycle)
	AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GatewayDestinationHandlerFunc)
	AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GatewayDestinationLifecycle)
}

type gatewayDestinationLister struct {
	controller *gatewayDestinationController
}

func (l *gatewayDestinationLister) List(namespace string, selector labels.Selector) (ret []*GatewayDestination, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*GatewayDestination))
	})
	return
}

func (l *gatewayDestinationLister) Get(namespace, name string) (*GatewayDestination, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    GatewayDestinationGroupVersionKind.Group,
			Resource: "gatewayDestination",
		}, key)
	}
	return obj.(*GatewayDestination), nil
}

type gatewayDestinationController struct {
	controller.GenericController
}

func (c *gatewayDestinationController) Generic() controller.GenericController {
	return c.GenericController
}

func (c *gatewayDestinationController) Lister() GatewayDestinationLister {
	return &gatewayDestinationLister{
		controller: c,
	}
}

func (c *gatewayDestinationController) AddHandler(ctx context.Context, name string, handler GatewayDestinationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GatewayDestination); ok {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

func (c *gatewayDestinationController) AddClusterScopedHandler(ctx context.Context, name, cluster string, handler GatewayDestinationHandlerFunc) {
	c.GenericController.AddHandler(ctx, name, func(key string, obj interface{}) (interface{}, error) {
		if obj == nil {
			return handler(key, nil)
		} else if v, ok := obj.(*GatewayDestination); ok && controller.ObjectInCluster(cluster, obj) {
			return handler(key, v)
		} else {
			return nil, nil
		}
	})
}

type gatewayDestinationFactory struct {
}

func (c gatewayDestinationFactory) Object() runtime.Object {
	return &GatewayDestination{}
}

func (c gatewayDestinationFactory) List() runtime.Object {
	return &GatewayDestinationList{}
}

func (s *gatewayDestinationClient) Controller() GatewayDestinationController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.gatewayDestinationControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(GatewayDestinationGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &gatewayDestinationController{
		GenericController: genericController,
	}

	s.client.gatewayDestinationControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type gatewayDestinationClient struct {
	client       *Client
	ns           string
	objectClient *objectclient.ObjectClient
	controller   GatewayDestinationController
}

func (s *gatewayDestinationClient) ObjectClient() *objectclient.ObjectClient {
	return s.objectClient
}

func (s *gatewayDestinationClient) Create(o *GatewayDestination) (*GatewayDestination, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*GatewayDestination), err
}

func (s *gatewayDestinationClient) Get(name string, opts metav1.GetOptions) (*GatewayDestination, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*GatewayDestination), err
}

func (s *gatewayDestinationClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*GatewayDestination, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*GatewayDestination), err
}

func (s *gatewayDestinationClient) Update(o *GatewayDestination) (*GatewayDestination, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*GatewayDestination), err
}

func (s *gatewayDestinationClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *gatewayDestinationClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *gatewayDestinationClient) List(opts metav1.ListOptions) (*GatewayDestinationList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*GatewayDestinationList), err
}

func (s *gatewayDestinationClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *gatewayDestinationClient) Patch(o *GatewayDestination, data []byte, subresources ...string) (*GatewayDestination, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*GatewayDestination), err
}

func (s *gatewayDestinationClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *gatewayDestinationClient) AddHandler(ctx context.Context, name string, sync GatewayDestinationHandlerFunc) {
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *gatewayDestinationClient) AddLifecycle(ctx context.Context, name string, lifecycle GatewayDestinationLifecycle) {
	sync := NewGatewayDestinationLifecycleAdapter(name, false, s, lifecycle)
	s.Controller().AddHandler(ctx, name, sync)
}

func (s *gatewayDestinationClient) AddClusterScopedHandler(ctx context.Context, name, clusterName string, sync GatewayDestinationHandlerFunc) {
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

func (s *gatewayDestinationClient) AddClusterScopedLifecycle(ctx context.Context, name, clusterName string, lifecycle GatewayDestinationLifecycle) {
	sync := NewGatewayDestinationLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.Controller().AddClusterScopedHandler(ctx, name, clusterName, sync)
}

type GatewayDestinationIndexer func(obj *GatewayDestination) ([]string, error)

type GatewayDestinationClientCache interface {
	Get(namespace, name string) (*GatewayDestination, error)
	List(namespace string, selector labels.Selector) ([]*GatewayDestination, error)

	Index(name string, indexer GatewayDestinationIndexer)
	GetIndexed(name, key string) ([]*GatewayDestination, error)
}

type GatewayDestinationClient interface {
	Create(*GatewayDestination) (*GatewayDestination, error)
	Get(namespace, name string, opts metav1.GetOptions) (*GatewayDestination, error)
	Update(*GatewayDestination) (*GatewayDestination, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	List(namespace string, opts metav1.ListOptions) (*GatewayDestinationList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)

	Cache() GatewayDestinationClientCache

	OnCreate(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc)
	OnChange(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc)
	OnRemove(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc)
	Enqueue(namespace, name string)

	Generic() controller.GenericController
	Interface() GatewayDestinationInterface
}

type gatewayDestinationClientCache struct {
	client *gatewayDestinationClient2
}

type gatewayDestinationClient2 struct {
	iface      GatewayDestinationInterface
	controller GatewayDestinationController
}

func (n *gatewayDestinationClient2) Interface() GatewayDestinationInterface {
	return n.iface
}

func (n *gatewayDestinationClient2) Generic() controller.GenericController {
	return n.iface.Controller().Generic()
}

func (n *gatewayDestinationClient2) Enqueue(namespace, name string) {
	n.iface.Controller().Enqueue(namespace, name)
}

func (n *gatewayDestinationClient2) Create(obj *GatewayDestination) (*GatewayDestination, error) {
	return n.iface.Create(obj)
}

func (n *gatewayDestinationClient2) Get(namespace, name string, opts metav1.GetOptions) (*GatewayDestination, error) {
	return n.iface.GetNamespaced(namespace, name, opts)
}

func (n *gatewayDestinationClient2) Update(obj *GatewayDestination) (*GatewayDestination, error) {
	return n.iface.Update(obj)
}

func (n *gatewayDestinationClient2) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return n.iface.DeleteNamespaced(namespace, name, options)
}

func (n *gatewayDestinationClient2) List(namespace string, opts metav1.ListOptions) (*GatewayDestinationList, error) {
	return n.iface.List(opts)
}

func (n *gatewayDestinationClient2) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return n.iface.Watch(opts)
}

func (n *gatewayDestinationClientCache) Get(namespace, name string) (*GatewayDestination, error) {
	return n.client.controller.Lister().Get(namespace, name)
}

func (n *gatewayDestinationClientCache) List(namespace string, selector labels.Selector) ([]*GatewayDestination, error) {
	return n.client.controller.Lister().List(namespace, selector)
}

func (n *gatewayDestinationClient2) Cache() GatewayDestinationClientCache {
	n.loadController()
	return &gatewayDestinationClientCache{
		client: n,
	}
}

func (n *gatewayDestinationClient2) OnCreate(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-create", &gatewayDestinationLifecycleDelegate{create: sync})
}

func (n *gatewayDestinationClient2) OnChange(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name+"-change", &gatewayDestinationLifecycleDelegate{update: sync})
}

func (n *gatewayDestinationClient2) OnRemove(ctx context.Context, name string, sync GatewayDestinationChangeHandlerFunc) {
	n.loadController()
	n.iface.AddLifecycle(ctx, name, &gatewayDestinationLifecycleDelegate{remove: sync})
}

func (n *gatewayDestinationClientCache) Index(name string, indexer GatewayDestinationIndexer) {
	err := n.client.controller.Informer().GetIndexer().AddIndexers(map[string]cache.IndexFunc{
		name: func(obj interface{}) ([]string, error) {
			if v, ok := obj.(*GatewayDestination); ok {
				return indexer(v)
			}
			return nil, nil
		},
	})

	if err != nil {
		panic(err)
	}
}

func (n *gatewayDestinationClientCache) GetIndexed(name, key string) ([]*GatewayDestination, error) {
	var result []*GatewayDestination
	objs, err := n.client.controller.Informer().GetIndexer().ByIndex(name, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if v, ok := obj.(*GatewayDestination); ok {
			result = append(result, v)
		}
	}

	return result, nil
}

func (n *gatewayDestinationClient2) loadController() {
	if n.controller == nil {
		n.controller = n.iface.Controller()
	}
}

type gatewayDestinationLifecycleDelegate struct {
	create GatewayDestinationChangeHandlerFunc
	update GatewayDestinationChangeHandlerFunc
	remove GatewayDestinationChangeHandlerFunc
}

func (n *gatewayDestinationLifecycleDelegate) HasCreate() bool {
	return n.create != nil
}

func (n *gatewayDestinationLifecycleDelegate) Create(obj *GatewayDestination) (runtime.Object, error) {
	if n.create == nil {
		return obj, nil
	}
	return n.create(obj)
}

func (n *gatewayDestinationLifecycleDelegate) HasFinalize() bool {
	return n.remove != nil
}

func (n *gatewayDestinationLifecycleDelegate) Remove(obj *GatewayDestination) (runtime.Object, error) {
	if n.remove == nil {
		return obj, nil
	}
	return n.remove(obj)
}

func (n *gatewayDestinationLifecycleDelegate) Updated(obj *GatewayDestination) (runtime.Object, error) {
	if n.update == nil {
		return obj, nil
	}
	return n.update(obj)
}
