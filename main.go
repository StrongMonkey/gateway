//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rancher/gateway/pkg/controllers/gateway"
	"github.com/rancher/gateway/pkg/server"
	appsv1 "github.com/rancher/gateway/types/apis/apps/v1beta2"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"github.com/rancher/norman"
	"github.com/rancher/norman/signal"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/apimachinery/pkg/util/proxy"
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

func run() error {
	logrus.Info("Starting controller")
	ctx := signal.SigTermCancelContext(context.Background())

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	normanConfig := server.Config()

	ctx, _, err = normanConfig.Build(ctx, &norman.Options{})
	if err != nil {
		return err
	}

	clients, err := v1.NewForConfig(*config)
	if err != nil {
		return err
	}
	cs := v1.NewClientsFromInterface(clients)
	gatewayHandler := Handler{
		gatewayDestLister: cs.GatewayDestination.Cache(),
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
	deploymentLister  appsv1.DeploymentClientCache
	deployments       appsv1.DeploymentInterface
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get(RioNameHeader)
	namespace := r.Header.Get(RioNamespaceHeader)
	gatewayDest, err := h.gatewayDestLister.Get(namespace, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	dep, err := h.deploymentLister.Get(gatewayDest.Spec.DestNamespace, gatewayDest.Spec.DestDeploymentName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	if dep.Spec.Replicas != nil && *dep.Spec.Replicas == 0 {
		dep.Spec.Replicas = &[]int32{1}[0]
		if _, err := h.deployments.Update(dep); err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
	}

	timer := time.After(time.Minute)

	ch, ok := gateway.EndpointChanMap.Load(fmt.Sprintf("%s.%s", name, namespace))
	if ok {
		c := ch.(chan struct{})
		select {
		case <-timer:
			http.Error(w, "timeout waiting for endpoint to be active", http.StatusGatewayTimeout)
		case _, ok := <-c:
			if !ok {
				targetUrl := &url.URL{
					Scheme: "http",
					Host:   fmt.Sprintf("%s.%s.svc.cluster.local", name, namespace),
				}
				port := r.URL.Port()
				if port != "" {
					targetUrl.Host = targetUrl.Host + ":" + port
				}
				r.URL = targetUrl
				r.Host = targetUrl.Host
				httpProxy := proxy.NewUpgradeAwareHandler(targetUrl, nil, false, false, er)
				httpProxy.ServeHTTP(w, r)
			}
		}
	}
	return
}

var er = &errorResponder{}

type errorResponder struct {
}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
