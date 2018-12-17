package v1

import (
	"github.com/rancher/norman/types"
	"github.com/rancher/norman/types/factory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	APIVersion = types.APIVersion{
		Group:   "some.api.group",
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
	MatchHeader        map[string]string
	MatchHost          string
	MatchPath          string
	DestServiceName    string
	DestDeploymentName string
}
