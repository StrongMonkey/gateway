package server

import (
	"github.com/rancher/gateway/pkg/controllers/gateway"
	gTypes "github.com/rancher/gateway/types"
	appsv1 "github.com/rancher/gateway/types/apis/apps/v1beta2"
	corev1 "github.com/rancher/gateway/types/apis/core/v1"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"github.com/rancher/gateway/types/client/gateway/v1"
	"github.com/rancher/norman"
	"github.com/rancher/norman/types"
)

func Config() *norman.Config {
	return &norman.Config{
		Name: "gateway",
		Schemas: []*types.Schemas{
			v1.Schemas,
		},

		CRDs: map[*types.APIVersion][]string{
			&v1.APIVersion: {
				client.GatewayDestinationType,
			},
		},

		Clients: []norman.ClientFactory{
			v1.Factory,
			appsv1.Factory,
			corev1.Factory,
		},

		MasterControllers: []norman.ControllerRegister{
			gTypes.Register(gateway.Register),
		},
	}
}
