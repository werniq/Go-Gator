package main

import (
	"flag"
	v1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"log"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

const (
	// tlsDir is a path to secret, which contains TLS certificates and key
	tlsDir = `/run/secrets/tls`

	// defaultTlsCertFile is the path to the certificate file used for serving the webhook over HTTPS
	defaultTlsCertFile = tlsDir + "/tls.crt"

	// defaultTlsKeyFile is the path to the private key file used for serving the webhook over HTTPS
	defaultTlsKeyFile = tlsDir + "/tls.key"
)

func init() {
	utilruntime.Must(v1.AddToScheme(scheme))

	utilruntime.Must(newsaggregatorv1.AddToScheme(scheme))
}

func main() {
	var (
		tlsCertFile string
		tlsKeyFile  string
	)
	flag.StringVar(&tlsCertFile, "c", defaultTlsCertFile,
		"Path to the certificate file used for serving the webhook over HTTPS")
	flag.StringVar(&tlsKeyFile, "k", defaultTlsKeyFile,
		"Path to the private key file used for serving the webhook over HTTPS")
	flag.Parsed()

	err := RunConfigMapController(tlsCertFile, tlsKeyFile)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("ConfigMap controller successfully started")
}
