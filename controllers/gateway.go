package controllers

import (
	"context"
	"fmt"
	"sync"

	"github.com/rancher/gateway/types"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type gatewayController struct {
	gateways v1.GatewayDestinationClient
}

func Register(ctx context.Context, rContext *types.Context) error {
	g := &gatewayController{
		gateways: rContext.Gateway.GatewayDestination,
	}
	rContext.Gateway.GatewayDestination.Interface().AddHandler(ctx, "gatewayDesitination", g.sync)
	return nil
}
