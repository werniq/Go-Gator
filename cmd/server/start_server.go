package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var (
	ErrRunningServer = fmt.Errorf("error while running server on port %s: ", Addr)
)

const (
	// Addr constant describes the port on which server will operate
	Addr = ":8080"
)

// ConfAndRun initializes server using gin framework, then attaches routes and handlers to it, and runs
// server on the port Addr
func ConfAndRun() {
	r := gin.Default()

	setupRoutes(r)

	go func() {
		err := FetchNewsJob()
		if err != nil {
			log.Fatalln("Error executing fetch news job: ", err)
		}
		time.Sleep(time.Hour * 24)
	}()

	err := r.Run(Addr)
	if err != nil {
		log.Fatalln(ErrRunningServer, err)
	}
}
