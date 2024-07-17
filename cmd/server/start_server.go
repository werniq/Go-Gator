package server

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogator/cmd/parsers"
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
)

const (
	// DevAddr constant describes the port on which server will operate in Development environment
	DevAddr = ":8080"

	// ProdAddr is used to run server in production environment
	ProdAddr = ":443"

	// CertFile is the name of certificate file
	CertFile = "certificate.pem"

	// KeyFile is the name of the key for the certificate above
	KeyFile = "key.pem"
)

// ConfAndRun initializes HTTPS server using gin framework, then attaches routes and handlers to it, and runs
// server on the port specified by user, or default - 443
func ConfAndRun() error {
	var (
		errChan = make(chan error, 1)
		server  = gin.Default()
		err     error

		// ServerPort identifies port on which Server will be running
		ServerPort int

		// UpdatesFrequency means every X hours after which new news will be parsed
		UpdatesFrequency int

		// certFile is the name of certificate file
		certFile string

		// keyFile is the name of the key for the certificate above
		keyFile string
	)
	flag.IntVar(&UpdatesFrequency, "f", 4, "How many hours fetch news job will wait after each execution")
	flag.IntVar(&ServerPort, "p", 443, "On which port server will be running")
	flag.StringVar(&certFile, "c", "certificate.pem", "Certificate for the HTTPs server")
	flag.StringVar(&keyFile, "k", "key.pem", "Private key for the HTTPs server")
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

	err = server.RunTLS(fmt.Sprintf(":%d", ServerPort),
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
