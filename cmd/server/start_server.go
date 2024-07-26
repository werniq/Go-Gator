package server

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	parsers "gogator/cmd/parsers"
	"gogator/cmd/server/handlers"
	"gogator/cmd/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	// defaultCertsPath is default path to server
	defaultCertsPath = filepath.Join("cmd", "server", "certs")
)

const (
	// defaultUpdateFrequency is an interval in hours of hours used to fetch and parse article feeds
	defaultUpdateFrequency = 4

	// defaultServerPort is a default port on which this server will be running
	defaultServerPort = 443

	// defaultCertName represents default name of server's certificate file
	defaultCertName = "certificate.pem"

	// defaultPrivateKey identifies the default name of server's private key
	defaultPrivateKey = "key.pem"

	// errRunFetchNews is thrown when we have problems while doing fetch news job
	errRunFetchNews = "error while doing fetch news job: "

	// errNotSpecified helps us to check if error was related to initializing sources file
	errNotSpecified = "The system cannot find the file specified."

	// errInitializingSources is thrown when func responsible for initialization of sources fails
	errInitializingSources = "Error initializing sources file:"
)

// ConfAndRun initializes and runs an HTTPS server using the Gin framework.
// This function sets up server routes and handlers, and starts the server
// on a user-specified port or defaults to port 443. It also launches a concurrent job
// which is fetching news feeds at a specified frequency.
//
// Optional parameters (specified via flags):
// / -f (updatesFrequency): Specifies the interval in hours at which the program
// /   will fetch and parse news feeds. Default value is used if not specified.
// / -p (serverPort): Specifies the port on which the server will run. Defaults to 443 if not specified.
// / -c (certFile): Specifies the absolute path to the certificate file for the HTTPS server. Defaults to a predefined path if not specified.
// / -k (keyFile): Specifies the absolute path to the private key file for the HTTPS server. Defaults to a predefined path if not specified.
// / -fs (storagePath): Specifies the path to the directory where all data will be stored. Defaults to a predefined path if not specified.
func ConfAndRun() error {
	var (
		errChan = make(chan error, 1)
		server  = gin.Default()
		err     error

		// serverPort identifies port on which Server will be running
		serverPort int

		// updatesFrequency means every X hours after which new news will be parsed
		updatesFrequency int

		// certFile is the name of certificate file
		certFile string

		// keyFile is the name of the key for the certificate above
		keyFile string

		storagePath string
	)
	cwdPath, err := os.Getwd()
	if err != nil {
		return err
	}

	flag.IntVar(&updatesFrequency, "f", defaultUpdateFrequency,
		"How many hours fetch news job will wait after each execution")
	flag.IntVar(&serverPort, "p", defaultServerPort,
		"On which port server will be running")
	flag.StringVar(&certFile, "c", filepath.Join(cwdPath, defaultCertsPath, defaultCertName),
		"Absolute path to the certificate for the HTTPs server")
	flag.StringVar(&keyFile, "k", filepath.Join(cwdPath, defaultCertsPath, defaultPrivateKey),
		"Absolute path to the private key for the HTTPs server")
	flag.StringVar(&storagePath, "fs", filepath.Join(parsers.CmdDir, parsers.ParsersDir, parsers.DataDir),
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

	go runFetchNewsJob(updatesFrequency, errChan)

	setupRoutes(server)

	go func() {
		errChan <- server.RunTLS(fmt.Sprintf(":%d", serverPort),
			certFile,
			keyFile)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	return nil
}

// runFetchNewsJob initializes and runs FetchNewsJob, which will parse data from feeds into respective files
func runFetchNewsJob(updatesFrequency int, errChan chan error) {
	dateTimestamp := time.Now().Format(time.DateOnly)
	j := FetchNewsJob{
		Filters: types.NewFilteringParams("", dateTimestamp, "", ""),
	}

	err := j.Run()
	if err != nil {
		errChan <- errors.New(errRunFetchNews + err.Error())
		return
	}

	time.Sleep(time.Hour * time.Duration(updatesFrequency))

	handlers.LastFetchedFileDate = time.Now().Format(time.DateOnly)
}
