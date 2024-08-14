/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"crypto/tls"
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	newsaggregatorv1 "teamdev.com/go-gator-operator/api/v1"
	"teamdev.com/go-gator-operator/internal/controller"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

const (
	// defaultMetricsBindAddress is the default address the metric endpoint should bind to
	defaultMetricsBindAddress = "0"

	// defaultHealthProbeBindAddress is the default address the health probe endpoint will bind to
	defaultHealthProbeBindAddress = ":8081"

	// defaultEnableLeaderElection is the default value for enabling leader election
	defaultEnableLeaderElection = false

	// defaultEnableSecureMetrics is the default value for enabling secure metrics
	defaultEnableSecureMetrics = true

	// defaultEnableHttp2 is the default value for enabling HTTP/2
	defaultEnableHttp2 = false

	// defaultServerUrl is a default url address of the news aggregator server
	defaultServerUrl = "https://go-gator-svc.go-gator.svc.cluster.local:443/news"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(newsaggregatorv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		probeAddr            string
		secureMetrics        bool
		serverUrl            string
		enableHTTP2          bool
		tlsOpts              []func(*tls.Config)
	)
	flag.StringVar(&metricsAddr, "metrics-bind-address", defaultMetricsBindAddress, "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	flag.StringVar(&serverUrl, "server-url", defaultServerUrl, "Url address of the news aggregator server")
	flag.StringVar(&probeAddr, "health-probe-bind-address", defaultHealthProbeBindAddress, "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", defaultEnableLeaderElection,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.BoolVar(&secureMetrics, "metrics-secure", defaultEnableSecureMetrics,
		"If set, the metrics endpoint is served securely via HTTPS. Use --metrics-secure=false to use HTTP instead.")
	flag.BoolVar(&enableHTTP2, "enable-http2", defaultEnableHttp2,
		"If set, HTTP/2 will be enabled for the metrics and webhook servers")

	opts := zap.Options{
		Development: true,
	}

	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}

	if !enableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}

	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})

	metricsServerOptions := metricsserver.Options{
		BindAddress:   metricsAddr,
		SecureServing: secureMetrics,
		TLSOpts:       tlsOpts,
	}

	if secureMetrics {
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		Metrics:                metricsServerOptions,
		WebhookServer:          webhookServer,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "dc074b3b.teamdev.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.HotNewsReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr, serverUrl); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HotNews")
		os.Exit(1)
	}
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&newsaggregatorv1.HotNews{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "HotNews")
			os.Exit(1)
		}
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
