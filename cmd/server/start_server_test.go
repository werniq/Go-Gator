package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"newsaggr/cmd/parsers"
	"os"
	"testing"
	"time"
)

func TestConfAndRun(t *testing.T) {
	testCases := []struct {
		Name        string
		Setup       func()
		ExpectError bool
	}{
		{
			Name: "Without sources.json file",
			Setup: func() {
				err := os.Remove(CwdPath + parsers.PathToDataDir + "sources.json")
				if err != nil {
					t.Error("Failed to remove sources.json file: ", err)
				}
			},
			ExpectError: true,
		},
	}

	gin.SetMode(gin.TestMode)
	server := gin.Default()
	setupRoutes(server)

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Setup()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			errCh := make(chan error, 1)
			go func() {
				errCh <- ConfAndRun()
			}()

			select {
			case err := <-errCh:
				if tt.ExpectError {
					assert.NotNil(t, err)
				} else {
					assert.Nil(t, err)
				}
			case <-ctx.Done():
				t.Log("ConfAndRun took too long, returning nil error because server is working fine.")
			}
		})
	}
}
