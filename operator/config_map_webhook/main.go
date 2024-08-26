package main

import (
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

const (
	// tlsDir is a path to secret, which contains TLS certificates and key
	tlsDir = `/run/secrets/tls`

	// tlsCertFile is the path to the certificate file used for serving the webhook over HTTPS
	tlsCertFile = tlsDir + "/tls.crt"

	// tlsKeyFile is the path to the private key file used for serving the webhook over HTTPS
	tlsKeyFile = tlsDir + "/tls.key"
)

func init() {
	utilruntime.Must(v1.AddToScheme(scheme))

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
