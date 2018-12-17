package main

import (
	"fmt"
	"path"

	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"github.com/rancher/norman/generator"
	"github.com/rancher/norman/types"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1beta2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	basePackage = "github.com/rancher/gateway/types"
	baseK8s     = "apis"
)

func main() {
	if err := generator.DefaultGenerate(v1.Schemas, "github.com/rancher/gateway/types", true, nil); err != nil {
		logrus.Fatal(err)
	}
	generateNativeTypes(corev1.SchemeGroupVersion, []interface{}{
		corev1.Endpoints{},
		corev1.Service{},
	}, nil)
	generateNativeTypes(appsv1.SchemeGroupVersion, []interface{}{
		appsv1.Deployment{},
	}, nil)
}

func generateNativeTypes(gv schema.GroupVersion, nsObjs []interface{}, objs []interface{}) {
	version := gv.Version
	group := gv.Group
	groupPath := group

	if groupPath == "" {
		groupPath = "core"
	}

	k8sOutputPackage := path.Join(basePackage, baseK8s, groupPath, version)

	if err := generator.GenerateControllerForTypes(&types.APIVersion{
		Version: version,
		Group:   group,
		Path:    fmt.Sprintf("/k8s/%s-%s", groupPath, version),
	}, k8sOutputPackage, nsObjs, objs); err != nil {
		panic(err)
	}
}
