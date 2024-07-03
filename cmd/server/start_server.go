package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"newsaggr/cmd/parsers"
	"newsaggr/cmd/server/handlers"
	"newsaggr/cmd/types"
	"os"
	"time"
)

var (
	// ErrRunningServer is throws whenever we encounter errors while running our server
	ErrRunningServer = fmt.Errorf("error while running server on port %s: ", DevAddr)

	// ErrFetchNewsJob is thrown when we have problems while doing fetch news job
	ErrFetchNewsJob = fmt.Errorf("error while doing fetch news job: ")
)

const (
	// DevAddr constant describes the port on which server will operate in Development environment
	DevAddr = ":8080"

	// ProdAddr is used to run server in production environment
	ProdAddr = ":443"

	// RelativePathToCertsDir is a path to the folder with OpenSSL Certificate and Key
	RelativePathToCertsDir = "\\cmd\\server\\certs\\"

	// CertFile is the name of certificate file
	CertFile = "server.crt"

	// KeyFile is the name of the key for the certificate above
	KeyFile = "server.key"

	// UpdatesFrequency is used to update the news every X hours
	UpdatesFrequency = 4

	CwdPath = "C:\\Users\\Oleksandr Matviienko\\GolandProjects\\newsAggr-2\\Go-Gator"
)

// ConfAndRun initializes server using gin framework, then attaches routes and handlers to it, and runs
// server on the port DevAddr
func ConfAndRun() error {
	var (
		errChan = make(chan error, 1)
		server  = gin.Default()
		err     error
	)

	err = parsers.LoadSourcesFile()
	switch {
	case errors.Is(err, os.ErrNotExist):
		err = parsers.InitSourcesFile()
		if err != nil {
			return err
		}
	case err != nil:
		return err
	}

	setupRoutes(server)

	go func() {
		dateTimestamp := time.Now().Format(time.DateOnly)
		j := FetchNewsJob{
			Filters: types.NewFilteringParams("", dateTimestamp, "", ""),
		}

		err := j.Run()
		if err != nil {
			errChan <- err
			return
		}

		time.Sleep(time.Hour * UpdatesFrequency)

		handlers.LastFetchedFileDate = time.Now().Format(time.DateOnly)
	}()

	//Cwd, err := os.Getwd()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//PathToCertsDir := Cwd + RelativePathToCertsDir
	//
	//err = server.RunTLS(ProdAddr, PathToCertsDir+CertFile, PathToCertsDir+KeyFile)
	//if err != nil {
	//	log.Fatalln(ErrRunningServer, err)
	//}

	err = server.Run(DevAddr)
	if err != nil {
		return err
	}

	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
