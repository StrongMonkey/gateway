//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/proxy"

	"github.com/rancher/gateway/controllers"
	gTypes "github.com/rancher/gateway/types"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"github.com/rancher/gateway/types/client/gateway/v1"
	"github.com/rancher/norman"
	"github.com/rancher/norman/signal"
	"github.com/rancher/norman/types"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

const (
	RioNameHeader      = "X-Rio-ServiceName"
	RioNamespaceHeader = "X-Rio-Namespace"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "gateway"
	app.Version = VERSION
	app.Usage = "gateway needs help!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "kubeconfig",
			EnvVar: "KUBECONFIG",
			Value:  "${HOME}/.kube/config",
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	logrus.Info("Starting controller")
	ctx := signal.SigTermCancelContext(context.Background())

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	normanConfig := &norman.Config{
		Name: "gateway",
		CRDs: map[*types.APIVersion][]string{
			&v1.APIVersion: {
				client.GatewayDestinationType,
			},
		},
		Clients: []norman.ClientFactory{
			v1.Factory,
		},
		Config:      config,
		GlobalSetup: gTypes.BuildContext,
		MasterControllers: []norman.ControllerRegister{
			gTypes.Register(controllers.Register),
		},
	}
	ctx, _, err = normanConfig.Build(ctx, &norman.Options{})
	if err != nil {
		return err
	}

	normanServer := norman.GetServer(ctx)
	clients, err := v1.NewForConfig(*config)
	if err != nil {
		return err
	}
	cs := v1.NewClientsFromInterface(clients)
	gatewayHandler := Handler{
		gatewayDestLister: cs.GatewayDestination.Cache(),
		appsV1:            normanServer.K8sClient.AppsV1(),
		corev1:            normanServer.K8sClient.CoreV1(),
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(gatewayHandler, &http2.Server{}),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("Error running HTTP server: %v", err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
	return nil
}

type Handler struct {
	gatewayDestLister v1.GatewayDestinationClientCache
	appsV1            appsv1.AppsV1Interface
	corev1            corev1.CoreV1Interface
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get(RioNameHeader)
	namespace := r.Header.Get(RioNamespaceHeader)
	gatewayDest, err := h.gatewayDestLister.Get(namespace, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	dep, err := h.appsV1.Deployments(gatewayDest.Spec.DestNamespace).Get(gatewayDest.Spec.DestDeploymentName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if dep.Spec.Replicas != nil && *dep.Spec.Replicas == 0 {
		dep.Spec.Replicas = &[]int32{1}[0]
		if _, err := h.appsV1.Deployments(gatewayDest.Spec.DestNamespace).Update(dep); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	// waiting for service endpoint > 0, then return FQDN
	service, err := h.corev1.Services(gatewayDest.Spec.DestNamespace).Get(gatewayDest.Spec.DestServiceName, metav1.GetOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	for i := 0; i < 15; i++ {
		endpoint, err := h.corev1.Endpoints(gatewayDest.Spec.DestNamespace).Get(gatewayDest.Spec.DestServiceName, metav1.GetOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		if len(endpoint.Subsets) > 0 {
			for _, port := range service.Spec.Ports {
				if strconv.Itoa(int(port.Port)) == r.URL.Port() {
					targetUrl := &url.URL{
						Scheme: "http",
						Host:   fmt.Sprintf("%s.%s.svc.cluster.local:%d", service.Name, service.Namespace, port.Port),
					}
					r.URL = targetUrl
					r.Host = targetUrl.Host
					httpProxy := proxy.NewUpgradeAwareHandler(targetUrl, nil, false, false, er)
					httpProxy.ServeHTTP(w, r)
					return
				}
			}
		}
		logrus.Debugf("Waiting for service %s to populate endpoints...", service.Name)
		time.Sleep(time.Second * 2)
		i++
	}

	http.Error(w, "timeout waiting for endpoint to be active", http.StatusGatewayTimeout)
	return
}

var er = &errorResponder{}

type errorResponder struct {
}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
