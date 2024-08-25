package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

var (
	scheme = runtime.NewScheme()
)

const (
	// tlsCertFile is the path to the certificate file used for serving the webhook over HTTPS
	tlsCertFile = "operator/config_map_webhook/tls/tls.crt"

	// tlsKeyFile is the path to the private key file used for serving the webhook over HTTPS
	tlsKeyFile = "operator/config_map_webhook/tls/tls.key"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(newsaggregatorv1.AddToScheme(scheme))
}

func main() {
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = RunConfigMapController(c)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("ConfigMap controller successfully started")
}
