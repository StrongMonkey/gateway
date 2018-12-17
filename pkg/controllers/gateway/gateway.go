package gateway

import (
	"context"
	"fmt"
	"sync"

	"github.com/rancher/gateway/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var EndpointChanMap = sync.Map{}

func Register(ctx context.Context, rContext *types.Context) error {
	h := handler{}
	rContext.Core.Endpoints.OnChange(ctx, "gateway-endpoint-watcher", h.Sync)
	rContext.Core.Endpoints.OnRemove(ctx, "gateway-endpoint-watcher", h.Remove)
	return nil
}

type handler struct{}

func (h handler) Sync(obj *corev1.Endpoints) (runtime.Object, error) {
	// todo: add a filter only for scale-to-zero services so that we don't have to keep a channel for every endpoint
	if obj != nil && obj.DeletionTimestamp == nil {
		ch := make(chan struct{}, 0)
		EndpointChanMap.LoadOrStore(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace), ch)
		if len(obj.Subsets) > 0 {
			o, _ := EndpointChanMap.Load(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace))
			c := o.(chan struct{})
			close(c)
		}
	}
	return obj, nil
}

func (h handler) Remove(obj *corev1.Endpoints) (runtime.Object, error) {
	if obj != nil {
		EndpointChanMap.Delete(fmt.Sprintf("%s.%s", obj.Name, obj.Namespace))
	}
	return obj, nil
}
