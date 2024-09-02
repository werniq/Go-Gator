package server

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	parsers "gogator/cmd/parsers"
	v1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"strings"
	"time"
)

var (
	// defaultCertsPath is default path to server
	defaultCertsPath = filepath.Join("cmd", "server", "certs")

	// defaultDataDirPath is a default path to the directory where all data will be stored
	defaultDataDirPath = filepath.Join("cmd", "parsers", "data")
	//defaultDataDirPath = "/tmp/"

	k8sClient client.Client
)

const (
	// defaultServerPort is a default port on which this server will be running
	defaultServerPort = 443

	// defaultCertificatesNamespace is a default namespace where certificates are stored
	defaultCertificatesNamespace = "go-gator"

	// defaultSecretName is a default name of the secret where certificates are stored
	defaultSecretName = "test-ca-secret"

	// defaultPrivateKey identifies the default name of server's private key
	defaultPrivateKey = "tls.key"

	// defaultCertName represents default name of server's certificate file
	defaultCertName = "tls.crt"

	// errNotSpecified helps us to check if error was related to initializing sources file
	errNotSpecified = "no such file or directory"

	// errInitializingSources is thrown when func responsible for initialization of sources fails
	errInitializingSources = "Error initializing sources file: "
)

// ConfAndRun initializes and runs an HTTPS server using the Gin framework.
// This function sets up server routes and handlers, and starts the server
// on a user-specified port or defaults to port 443.
//
// Optional parameters (specified via flags):
// / -p (serverPort): Specifies the port on which the server will run. Defaults to 443 if not specified.
// / -c (certFile): Specifies the absolute path to the certificate file for the HTTPS server. Defaults to a predefined path if not specified.
// / -k (keyFile): Specifies the absolute path to the private key file for the HTTPS server. Defaults to a predefined path if not specified.
// / -fs (storagePath): Specifies the path to the directory where all data will be stored. Defaults to a predefined path if not specified.
func ConfAndRun() error {
	var (
		server = gin.Default()
		err    error

		// serverPort identifies port on which Server will be running
		serverPort int

		// certFile is the name of certificate file
		certFile string

		// keyFile is the name of the key for the certificate above
		keyFile string

		// storagePath is a path where all data from application will be stored (sources and files with articles)
		storagePath string
	)
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	flag.IntVar(&serverPort, "p", defaultServerPort,
		"On which port server will be running")
	flag.StringVar(&certFile, "c", filepath.Join(cwdPath, defaultCertsPath, defaultCertName),
		"Absolute path to the certificate for the HTTPs server")
	flag.StringVar(&keyFile, "k", filepath.Join(cwdPath, defaultCertsPath, defaultPrivateKey),
		"Absolute path to the private key for the HTTPs server")
	flag.StringVar(&storagePath, "fs", defaultDataDirPath,
		"Path to directory where all data will be stored")
	flag.Parse()

	c := config.GetConfigOrDie()

	k8sClient, err = client.New(c, client.Options{})
	if err != nil {
		return err
	}
	parsers.StoragePath = storagePath

	err = parsers.LoadSourcesFile()
	if err != nil {
		if strings.Contains(err.Error(), errNotSpecified) {
			err = parsers.UpdateSourceFile()
			if err != nil {
				return errors.New(errInitializingSources + err.Error())
			}
		} else {
			return err
		}
	}

	setupRoutes(server)

	certPath, keyPath, err := loadCertsFromSecrets()

	err = server.RunTLS(fmt.Sprintf(":%d", serverPort),
		certPath,
		keyPath)
	if err != nil {
		return err
	}

	return nil
}

// loadCertsFromSecrets loads certificates from Kubernetes secrets
func loadCertsFromSecrets() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var k8sSecret v1.Secret
	err := k8sClient.Get(ctx, client.ObjectKey{
		Name: defaultSecretName,
	}, &k8sSecret)
	if err != nil {
		return "", "", err
	}

	certData := k8sSecret.Data[defaultCertName]
	keyData := k8sSecret.Data[defaultPrivateKey]

	cwdPath, err := os.Getwd()
	if err != nil {
		return "", "", err

	}

	defaultCertPath := filepath.Join(cwdPath, defaultCertsPath, defaultCertName)
	err = createFileFromDataAndPath(certData, defaultCertPath)
	if err != nil {
		return "", "", err
	}

	defaultPrivateKeyPath := filepath.Join(cwdPath, defaultCertsPath, defaultPrivateKey)
	err = createFileFromDataAndPath(keyData, defaultPrivateKeyPath)
	if err != nil {
		return "", "", err
	}

	return defaultCertPath, defaultPrivateKeyPath, nil
}

// createFileFromDataAndPath creates a file based on given file data and path
func createFileFromDataAndPath(fileData []byte, filepath string) error {
	err := os.WriteFile(filepath, fileData, 0644)
	if err != nil {
		return err
	}

	return nil
}
