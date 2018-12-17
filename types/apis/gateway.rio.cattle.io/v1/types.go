package v1

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/factory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	APIVersion = types.APIVersion{
		Group:   "gateway.rio.cattle.io",
		Version: "v1",
		Path:    "/v1-someapi",
	}
	Schemas = factory.
		Schemas(&APIVersion).
		MustImport(&APIVersion, GatewayDestination{})
)

type GatewayDestination struct {
	types.Namespaced

	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GatewayDestinationSpec `json:"spec"`
}

type GatewayDestinationSpec struct {
	MatchHeader        map[string]string `json:"matchHeader,omitempty"`
	MatchHost          string            `json:"matchHost,omitempty"`
	MatchPath          string            `json:"matchPath,omitempty"`
	DestServiceName    string            `json:"destServiceName,omitempty"`
	DestNamespace      string            `json:"destNamespace,omitempty"`
	DestDeploymentName string            `json:"destDeploymentName,omitempty"`
}
