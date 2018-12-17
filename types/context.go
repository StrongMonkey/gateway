package types

import (
	"context"

	appsv1 "github.com/rancher/gateway/types/apis/apps/v1beta2"
	corev1 "github.com/rancher/gateway/types/apis/core/v1"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
)

type contextKey struct{}

type Context struct {
	Apps    *appsv1.Clients
	Core    *corev1.Clients
	Gateway *v1.Clients
}

func Store(ctx context.Context, c *Context) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

func From(ctx context.Context) *Context {
	return ctx.Value(contextKey{}).(*Context)
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Gateway: v1.ClientsFrom(ctx),
		Apps:    appsv1.ClientsFrom(ctx),
		Core:    corev1.ClientsFrom(ctx),
	}
}

func BuildContext(ctx context.Context) (context.Context, error) {
	return Store(ctx, NewContext(ctx)), nil
}

func Register(f func(context.Context, *Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		return f(ctx, From(ctx))
	}
}
