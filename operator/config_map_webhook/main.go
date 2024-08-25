package main

import (
	"crypto/tls"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"log"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

var (
	scheme = runtime.NewScheme()
)

const (
	tlsDir = `/run/secrets/tls`

	// tlsCertFile is the path to the certificate file used for serving the webhook over HTTPS
	tlsCertFile = tlsDir + "/tls.crt"

	// tlsKeyFile is the path to the private key file used for serving the webhook over HTTPS
	tlsKeyFile = tlsDir + "tls.key"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(newsaggregatorv1.AddToScheme(scheme))
}

func main() {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}
	c, err := client.New(ctrl.GetConfigOrDie(), client.Options{
		HTTPClient:     customClient,
		Scheme:         scheme,
		Mapper:         nil,
		Cache:          nil,
		WarningHandler: client.WarningHandlerOptions{},
		DryRun:         nil,
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
