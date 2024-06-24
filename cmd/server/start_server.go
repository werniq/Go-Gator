package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var (
	ErrRunningServer = fmt.Errorf("error while running server on port %s: ", DevAddr)
)

const (
	// DevAddr constant describes the port on which server will operate in Development environment
	DevAddr = ":8080"

	// ProdAddr is used to run server in production environment
	ProdAddr = ":443"

	RelativePathToCertsDir = "\\cmd\\server\\certs\\"

	// CertFile is the name of certificate file
	CertFile = "server.crt"

	// KeyFile is the name of the key for the certificate above
	KeyFile = "server.key"

	// UpdatesFrequency is used to update the news every X hours
	UpdatesFrequency = 4
)

// ConfAndRun initializes server using gin framework, then attaches routes and handlers to it, and runs
// server on the port DevAddr
func ConfAndRun() {
	r := gin.Default()

	setupRoutes(r)

	go func() {
		err := FetchNewsJob()
		if err != nil {
			log.Fatalln("Error executing fetch news job: ", err)
		}
		time.Sleep(time.Hour * UpdatesFrequency)
	}()

	//Cwd, err := os.Getwd()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//PathToCertsDir := Cwd + RelativePathToCertsDir
	//
	//err = r.RunTLS(ProdAddr, PathToCertsDir+CertFile, PathToCertsDir+KeyFile)
	//if err != nil {
	//	log.Fatalln(ErrRunningServer, err)
	//}

	err := r.Run(DevAddr)
	if err != nil {
		log.Fatalln(ErrRunningServer, err)
	}
}
