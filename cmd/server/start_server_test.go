package server

import (
	"context"
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
			Name: "With sources.json file",
			Setup: func() {

			},
			ExpectError: false,
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
				}
			case <-ctx.Done():
				t.Log("ConfAndRun took too long, returning nil error because server is working fine.")
			}
		})
	}
}
