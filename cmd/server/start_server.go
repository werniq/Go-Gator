package server

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	parsers "gogator/cmd/parsers"
	"gogator/cmd/server/handlers"
	"gogator/cmd/types"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	// ErrFetchNewsJob is thrown when we have problems while doing fetch news job
	ErrFetchNewsJob = fmt.Errorf("error while doing fetch news job: ")

	// ErrRunningServer is thrown when we have error while running
	ErrRunningServer = fmt.Errorf("error running server: ")

	// RelativePathToCertsDir is a path to the folder with OpenSSL Certificate and Key
	RelativePathToCertsDir = filepath.Join("cmd", "server", "certs")

	DefaultCertPaths = filepath.Join("cmd", "server", "certs")
)

const (
	DefaultUpdatesFrequency = 4

	DefaultServerPort = 443

	DefaultCertName = "certificate.pem"

	DefaultPkey = "key.pem"
)

// ConfAndRun initializes HTTPS server using gin framework, then attaches routes and handlers to it, and runs
// server on the port specified by user, or default - 443
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
		log.Fatalln(err)
	}

	flag.IntVar(&updatesFrequency, "f", DefaultUpdatesFrequency,
		"How many hours fetch news job will wait after each execution")
	flag.IntVar(&serverPort, "p", DefaultServerPort,
		"On which port server will be running")
	flag.StringVar(&certFile, "c", filepath.Join(cwdPath, DefaultCertPaths, DefaultCertName),
		"Absolute path to the certificate for the HTTPs server")
	flag.StringVar(&keyFile, "k", filepath.Join(cwdPath, DefaultCertPaths, DefaultPkey),
		"Absolute path to the private key for the HTTPs server")
	flag.StringVar(&storagePath, "fs", filepath.Join(parsers.CmdDir, parsers.ParsersDir, parsers.DataDir),
		"Path to directory where all data will be stored")
	flag.Parse()

	parsers.StoragePath = storagePath

	err = parsers.LoadSourcesFile()
	if err != nil {
		err = parsers.UpdateSourcesFile()
		if err != nil {
			log.Println("Error initializing sources file: ", err.Error())
			return err
		}
	}

	go runFetchNewsJob(updatesFrequency, errChan)

	setupRoutes(server)

	err = server.RunTLS(fmt.Sprintf(":%d", serverPort),
		certFile,
		keyFile)

	if err != nil {
		log.Fatalln(ErrRunningServer, err)
	}

	close(errChan)

	for err := range errChan {
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
		log.Println(ErrFetchNewsJob, err.Error())
		errChan <- err
		return
	}

	time.Sleep(time.Hour * time.Duration(updatesFrequency))

	handlers.LastFetchedFileDate = time.Now().Format(time.DateOnly)
}
