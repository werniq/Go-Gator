package server

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/server/handlers"
	"newsaggr/cmd/types"
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
)

const (
	// prodAddr is used to run server in production environment
	prodAddr = ":443"
)

// ConfAndRun initializes HTTPS server using gin framework, then attaches routes and handlers to it, and runs
// server on the port prodAddr
func ConfAndRun() error {
	var (
		errChan = make(chan error, 1)
		server  = gin.Default()
		err     error

		// UpdatesFrequency means every X hours after which new news will be parsed
		UpdatesFrequency int

		// certFile is the name of certificate file
		certFile string

		// keyFile is the name of the key for the certificate above
		keyFile string
	)
	flag.IntVar(&UpdatesFrequency, "f", 4, "How many hours fetch news job will wait after each execution")
	flag.StringVar(&certFile, "cert", "certificate.pem", "Certificate for the HTTPs server")
	flag.StringVar(&keyFile, "key", "key.pem", "Private key for the HTTPs server")
	flag.Parse()

	err = parsers.LoadSourcesFile()
	if err != nil {
		err = parsers.InitSourcesFile()
		if err != nil {
			log.Println("Error initializing sources file: ", err.Error())
			return err
		}
	}

	go func() {
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

		time.Sleep(time.Hour * time.Duration(UpdatesFrequency))

		handlers.LastFetchedFileDate = time.Now().Format(time.DateOnly)
	}()

	Cwd, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	PathToCertsDir := filepath.Join(Cwd, RelativePathToCertsDir)

	setupRoutes(server)

	err = server.RunTLS(prodAddr,
		filepath.Join(PathToCertsDir, certFile),
		filepath.Join(PathToCertsDir, keyFile))
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
