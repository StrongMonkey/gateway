//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"context"
	"fmt"
	"github.com/knative/pkg/logging"
	"github.com/knative/pkg/logging/logkey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/url"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
	activatorutil "github.com/knative/serving/pkg/activator/util"
	"github.com/rancher/gateway/pkg/controllers/gateway"
	"github.com/rancher/gateway/pkg/server"
	"github.com/rancher/gateway/types"
	appsv1 "github.com/rancher/gateway/types/apis/apps/v1beta2"
	"github.com/rancher/gateway/types/apis/gateway.rio.cattle.io/v1"
	"github.com/rancher/norman"
	"github.com/rancher/norman/signal"
	riov1 "github.com/rancher/rio/types/apis/rio.cattle.io/v1"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"k8s.io/apimachinery/pkg/util/proxy"
)

const (
	RioNameHeader          = "X-Rio-ServiceName"
	RioNamespaceHeader     = "X-Rio-Namespace"
	maxRetries             = 18 // the sum of all retries would add up to 1 minute
	minRetryInterval       = 100 * time.Millisecond
	exponentialBackoffBase = 1.3
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

	normanConfig := server.Config()

	ctx, _, err := normanConfig.Build(ctx, &norman.Options{})
	if err != nil {
		return err
	}

	rContext := types.From(ctx)
	gatewayHandler := Handler{
		services:          rContext.Rio.Service,
		serviceLister:     rContext.Rio.Service.Cache(),
		gatewayDestLister: rContext.Gateway.GatewayDestination.Cache(),
		deploymentLister:  rContext.Apps.Deployment.Cache(),
		deployments:       rContext.Apps.Deployment,
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(gatewayHandler, &http2.Server{}),
	}

	go func() {
		logrus.Infof("starting gateway server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.Errorf("Error running HTTP server: %v", err)
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
	return nil
}

type Handler struct {
	serviceLister     riov1.ServiceClientCache
	services          riov1.ServiceClient
	gatewayDestLister v1.GatewayDestinationClientCache
	deploymentLister  appsv1.DeploymentClientCache
	deployments       appsv1.DeploymentClient
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get(RioNameHeader)
	namespace := r.Header.Get(RioNamespaceHeader)
	gatewayDest, err := h.gatewayDestLister.Get(namespace, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	if os.Getenv("RIO_IN_CLUSTER") != "" {
		rioSvc, err := h.serviceLister.Get(namespace, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		if rioSvc.Spec.Scale == 0 {
			rioSvc.Spec.Scale = 1
			if _, err := h.services.Update(rioSvc); err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
		}
	} else {
		dep, err := h.deploymentLister.Get(gatewayDest.Spec.DestNamespace, gatewayDest.Spec.DestDeploymentName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		if dep.Spec.Replicas != nil || *dep.Spec.Replicas == 0 {
			dep.Spec.Replicas = &[]int32{1}[0]
			if _, err := h.deployments.Update(dep); err != nil {
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
				return
			}
		}
	}

	timer := time.After(time.Minute)

	endpointCh, endpointNotReady := gateway.EndpointChanMap.Load(fmt.Sprintf("%s.%s", name, namespace))
	if !endpointNotReady {
		serveFQDN(name, namespace, w, r)
		return
	}
	select {
	case <-timer:
		http.Error(w, "timeout waiting for endpoint to be active", http.StatusGatewayTimeout)
		return
	case _, ok := <-endpointCh.(chan struct{}):
		if !ok {
			serveFQDN(name, namespace, w, r)
			return
		}
	}
}

func serveFQDN(name, namespace string, w http.ResponseWriter, r *http.Request) {
	targetUrl := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s.%s.svc.cluster.local", name, namespace),
		Path:   r.URL.Path,
	}
	port := r.URL.Port()
	if port != "" {
		targetUrl.Host = targetUrl.Host + ":" + port
	}
	r.URL = targetUrl
	r.Host = targetUrl.Host

	// todo: check if 503 is actually coming from application or envoy
	shouldRetry := activatorutil.RetryStatus(http.StatusServiceUnavailable)
	backoffSettings := wait.Backoff{
		Duration: minRetryInterval,
		Factor:   exponentialBackoffBase,
		Steps:    maxRetries,
	}

	createdLogger, _ := logging.NewLogger("", zapcore.InfoLevel.String())
	logger := createdLogger.With(zap.String(logkey.ControllerType, "rio-autoscaler-gateway"))
	defer logger.Sync()

	rt := activatorutil.NewRetryRoundTripper(activatorutil.AutoTransport, logger, backoffSettings, shouldRetry)
	httpProxy := proxy.NewUpgradeAwareHandler(targetUrl, rt, true, false, er)
	httpProxy.ServeHTTP(w, r)
}

var er = &errorResponder{}

type errorResponder struct {
}

func (e *errorResponder) Error(w http.ResponseWriter, req *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
