package client

import (
	"github.com/rancher/norman/types"
)

const (
	GatewayDestinationType                    = "gatewayDestination"
	GatewayDestinationFieldAnnotations        = "annotations"
	GatewayDestinationFieldCreated            = "created"
	GatewayDestinationFieldDestDeploymentName = "destDeploymentName"
	GatewayDestinationFieldDestNamespace      = "destNamespace"
	GatewayDestinationFieldDestServiceName    = "destServiceName"
	GatewayDestinationFieldLabels             = "labels"
	GatewayDestinationFieldMatchHeader        = "matchHeader"
	GatewayDestinationFieldMatchHost          = "matchHost"
	GatewayDestinationFieldMatchPath          = "matchPath"
	GatewayDestinationFieldName               = "name"
	GatewayDestinationFieldNamespace          = "namespace"
	GatewayDestinationFieldOwnerReferences    = "ownerReferences"
	GatewayDestinationFieldRemoved            = "removed"
	GatewayDestinationFieldUUID               = "uuid"
)

type GatewayDestination struct {
	types.Resource
	Annotations        map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	Created            string            `json:"created,omitempty" yaml:"created,omitempty"`
	DestDeploymentName string            `json:"destDeploymentName,omitempty" yaml:"destDeploymentName,omitempty"`
	DestNamespace      string            `json:"destNamespace,omitempty" yaml:"destNamespace,omitempty"`
	DestServiceName    string            `json:"destServiceName,omitempty" yaml:"destServiceName,omitempty"`
	Labels             map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	MatchHeader        map[string]string `json:"matchHeader,omitempty" yaml:"matchHeader,omitempty"`
	MatchHost          string            `json:"matchHost,omitempty" yaml:"matchHost,omitempty"`
	MatchPath          string            `json:"matchPath,omitempty" yaml:"matchPath,omitempty"`
	Name               string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace          string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	OwnerReferences    []OwnerReference  `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	Removed            string            `json:"removed,omitempty" yaml:"removed,omitempty"`
	UUID               string            `json:"uuid,omitempty" yaml:"uuid,omitempty"`
}

type GatewayDestinationCollection struct {
	types.Collection
	Data   []GatewayDestination `json:"data,omitempty"`
	client *GatewayDestinationClient
}

type GatewayDestinationClient struct {
	apiClient *Client
}

type GatewayDestinationOperations interface {
	List(opts *types.ListOpts) (*GatewayDestinationCollection, error)
	Create(opts *GatewayDestination) (*GatewayDestination, error)
	Update(existing *GatewayDestination, updates interface{}) (*GatewayDestination, error)
	Replace(existing *GatewayDestination) (*GatewayDestination, error)
	ByID(id string) (*GatewayDestination, error)
	Delete(container *GatewayDestination) error
}

func newGatewayDestinationClient(apiClient *Client) *GatewayDestinationClient {
	return &GatewayDestinationClient{
		apiClient: apiClient,
	}
}

func (c *GatewayDestinationClient) Create(container *GatewayDestination) (*GatewayDestination, error) {
	resp := &GatewayDestination{}
	err := c.apiClient.Ops.DoCreate(GatewayDestinationType, container, resp)
	return resp, err
}

func (c *GatewayDestinationClient) Update(existing *GatewayDestination, updates interface{}) (*GatewayDestination, error) {
	resp := &GatewayDestination{}
	err := c.apiClient.Ops.DoUpdate(GatewayDestinationType, &existing.Resource, updates, resp)
	return resp, err
}

func (c *GatewayDestinationClient) Replace(obj *GatewayDestination) (*GatewayDestination, error) {
	resp := &GatewayDestination{}
	err := c.apiClient.Ops.DoReplace(GatewayDestinationType, &obj.Resource, obj, resp)
	return resp, err
}

func (c *GatewayDestinationClient) List(opts *types.ListOpts) (*GatewayDestinationCollection, error) {
	resp := &GatewayDestinationCollection{}
	err := c.apiClient.Ops.DoList(GatewayDestinationType, opts, resp)
	resp.client = c
	return resp, err
}

func (cc *GatewayDestinationCollection) Next() (*GatewayDestinationCollection, error) {
	if cc != nil && cc.Pagination != nil && cc.Pagination.Next != "" {
		resp := &GatewayDestinationCollection{}
		err := cc.client.apiClient.Ops.DoNext(cc.Pagination.Next, resp)
		resp.client = cc.client
		return resp, err
	}
	return nil, nil
}

func (c *GatewayDestinationClient) ByID(id string) (*GatewayDestination, error) {
	resp := &GatewayDestination{}
	err := c.apiClient.Ops.DoByID(GatewayDestinationType, id, resp)
	return resp, err
}

func (c *GatewayDestinationClient) Delete(container *GatewayDestination) error {
	return c.apiClient.Ops.DoResourceDelete(GatewayDestinationType, &container.Resource)
}
