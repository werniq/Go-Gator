package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/server/handlers"
	"newsaggr/cmd/types"
	"os"
	"path/filepath"
	"strconv"
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

	// ProdAddr is used to run server in production environment
	ProdAddr = ":443"

	// CertFile is the name of certificate file
	CertFile = "certificate.pem"

	// KeyFile is the name of the key for the certificate above
	KeyFile = "key.pem"
)

// ConfAndRun initializes HTTPS server using gin framework, then attaches routes and handlers to it, and runs
// server on the port ProdAddr
func ConfAndRun() error {
	var (
		errChan             = make(chan error, 1)
		server              = gin.Default()
		err                 error
		updatesFrequencyStr = os.Getenv("FETCH_NEWS_UPDATES_FREQUENCY")
		UpdatesFrequency    int
	)
	if updatesFrequencyStr != "" {
		UpdatesFrequency, err = strconv.Atoi(updatesFrequencyStr)
		if err != nil {
			return err
		}
	} else {
		UpdatesFrequency = 4
	}

	err = parsers.LoadSourcesFile()
	if err != nil {
		err = parsers.InitSourcesFile()
		if err != nil {
			log.Println("Error initializing sources file: ", err.Error())
			return err
		}
	}

	setupRoutes(server)
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

	err = server.RunTLS(ProdAddr,
		filepath.Join(PathToCertsDir, CertFile),
		filepath.Join(PathToCertsDir, KeyFile))
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
