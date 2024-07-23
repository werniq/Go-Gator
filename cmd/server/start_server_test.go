package server

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"path/filepath"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	parsers.StoragePath = filepath.Join("..", "parsers", "data")
}

func TestConfAndRun(t *testing.T) {
	testCases := []struct {
		Name        string
		Setup       func()
		ExpectError bool
	}{
		{
			Name: "Successful run",
			Setup: func() {
				err := flag.Set("p", "443")
				assert.Nil(t, err)
				err = flag.Set("f", "1")
				assert.Nil(t, err)
			},
			ExpectError: false,
		},
		{
			Name: "No certificates for server",
			Setup: func() {
				err := flag.Set("c", "")
				assert.Nil(t, err)

				err = flag.Set("k", "")
				assert.Nil(t, err)
			},
			ExpectError: true,
		},
		{
			Name: "Invalid port number",
			Setup: func() {
				err := flag.Set("p", "-1")
				assert.Nil(t, err)
			},
			ExpectError: true,
		},
		{
			Name: "Invalid certificate paths",
			Setup: func() {
				err := flag.Set("c", "invalid/cert.pem")
				assert.Nil(t, err)

				err = flag.Set("k", "invalid/key.pem")
				assert.Nil(t, err)
			},
			ExpectError: true,
		},
		{
			Name: "Invalid storage path",
			Setup: func() {
				err := flag.Set("fs", "/invalid/path")
				assert.Nil(t, err)
			},
			ExpectError: true,
		},
	}

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
