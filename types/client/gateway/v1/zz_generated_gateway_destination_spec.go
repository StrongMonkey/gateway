package client

const (
	GatewayDestinationSpecType                    = "gatewayDestinationSpec"
	GatewayDestinationSpecFieldDestDeploymentName = "destDeploymentName"
	GatewayDestinationSpecFieldDestNamespace      = "destNamespace"
	GatewayDestinationSpecFieldDestServiceName    = "destServiceName"
	GatewayDestinationSpecFieldMatchHeader        = "matchHeader"
	GatewayDestinationSpecFieldMatchHost          = "matchHost"
	GatewayDestinationSpecFieldMatchPath          = "matchPath"
)

type GatewayDestinationSpec struct {
	DestDeploymentName string            `json:"destDeploymentName,omitempty" yaml:"destDeploymentName,omitempty"`
	DestNamespace      string            `json:"destNamespace,omitempty" yaml:"destNamespace,omitempty"`
	DestServiceName    string            `json:"destServiceName,omitempty" yaml:"destServiceName,omitempty"`
	MatchHeader        map[string]string `json:"matchHeader,omitempty" yaml:"matchHeader,omitempty"`
	MatchHost          string            `json:"matchHost,omitempty" yaml:"matchHost,omitempty"`
	MatchPath          string            `json:"matchPath,omitempty" yaml:"matchPath,omitempty"`
}
