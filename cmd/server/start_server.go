package server

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	parsers "gogator/cmd/parsers"
	"os"
	"path/filepath"
	"strings"
)

var (
	// defaultCertsPath is default path to server
	defaultCertsPath = filepath.Join("cmd", "server", "certs")

	defaultDataDirPath = "/tmp/"
)

const (
	// defaultServerPort is a default port on which this server will be running
	defaultServerPort = 443

	// defaultCertName represents default name of server's certificate file
	defaultCertName = "certificate.pem"

	// defaultPrivateKey identifies the default name of server's private key
	defaultPrivateKey = "key.pem"

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

	err = server.RunTLS(fmt.Sprintf(":%d", serverPort),
		certFile,
		keyFile)
	if err != nil {
		return err
	}

	return nil
}
