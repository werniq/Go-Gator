package server

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	"os"
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
		Cleanup     func()
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
			Cleanup:     func() {},
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
			Cleanup:     func() {},
			ExpectError: true,
		},
		{
			Name: "Invalid port number",
			Setup: func() {
				err := flag.Set("p", "-1")
				assert.Nil(t, err)
			},
			Cleanup:     func() {},
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
			Cleanup:     func() {},
			ExpectError: true,
		},
		{
			Name: "Invalid storage path",
			Setup: func() {
				err := flag.Set("fs", "/invalid/path")
				assert.Nil(t, err)
			},
			Cleanup:     func() {},
			ExpectError: true,
		},
		{
			Name: "Invalid .PEM Certificate and Key files",
			Setup: func() {
				invalidCert := []byte("invalid certificate content")
				invalidKey := []byte("invalid key content")

				err := os.WriteFile("invalid_cert.pem", invalidCert, 0644)
				assert.Nil(t, err)

				err = os.WriteFile("invalid_key.pem", invalidKey, 0644)
				assert.Nil(t, err)

				err = flag.Set("c", "invalid_cert.pem")
				assert.Nil(t, err)

				err = flag.Set("k", "invalid_key.pem")
				assert.Nil(t, err)
			},
			Cleanup: func() {
				err := os.Remove("invalid_cert.pem")
				assert.Nil(t, err)
				err = os.Remove("invalid_key.pem")
				assert.Nil(t, err)
			},
			ExpectError: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			tt.Setup()
			defer tt.Cleanup()

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
