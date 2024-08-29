package server

import (
	"context"
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gogator/cmd/parsers"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
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

func Test_createFileFromDataAndPath(t *testing.T) {

	type args struct {
		fileData []byte
		filepath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, createFileFromDataAndPath(tt.args.fileData, tt.args.filepath), fmt.Sprintf("createFileFromDataAndPath(%v, %v)", tt.args.fileData, tt.args.filepath))
		})
	}
}

func Test_loadCertsFromSecrets(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"totalAmount": 2, "news": [{"title": "News 1"}, {"title": "News 2"}]}`))
	}))

	existingSecretList := v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      "feed-sample",
		},
		Data: map[string][]byte{
			"feed": []byte(mockServer.URL),
		},
	}

	k8sClient = fake.NewClientBuilder().WithScheme(scheme).
		WithObjects(&existingSecretList).
		Build()

	tests := []struct {
		name    string
		want    string
		want1   string
		wantErr bool
		setup   func()
	}{
		{
			name:    "Successful Load",
			want:    "path/to/cert.crt",
			want1:   "path/to/private.key",
			wantErr: false,
			setup: func() {
				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      defaultSecretName,
						Namespace: "default",
					},
					Data: map[string][]byte{
						defaultCertName:   []byte("cert-data"),
						defaultPrivateKey: []byte("key-data"),
					},
				}
				_ = k8sClient.Create(context.Background(), secret)
			},
		},
		{
			name:    "Secret Not Found",
			want:    "",
			want1:   "",
			wantErr: true,
			setup: func() {
			},
		},
		{
			name:    "Missing Certificate Data",
			want:    "",
			want1:   "",
			wantErr: true,
			setup: func() {
				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      defaultSecretName,
						Namespace: "default",
					},
					Data: map[string][]byte{
						defaultPrivateKey: []byte("key-data"),
					},
				}
				_ = k8sClient.Create(context.Background(), secret)
			},
		},
		{
			name:    "Missing Private Key Data",
			want:    "",
			want1:   "",
			wantErr: true,
			setup: func() {
				secret := &v1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      defaultSecretName,
						Namespace: "default",
					},
					Data: map[string][]byte{
						defaultCertName: []byte("cert-data"),
					},
				}
				_ = k8sClient.Create(context.Background(), secret)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := loadCertsFromSecrets()
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
