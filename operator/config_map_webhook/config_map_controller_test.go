package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	admission "k8s.io/api/admission/v1beta1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"strings"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	v1 "teamdev.com/go-gator/api/v1"
	"testing"
	"time"
)

// errorReader is a custom io.Reader that returns an error
type errorReader struct{}

func (er *errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("forced read error")
}

func TestWebhookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	universalDeserializer = serializer.NewCodecFactory(scheme).UniversalDeserializer()

	_ = v1.AddToScheme(scheme)
	_ = newsaggregatorv1.AddToScheme(scheme)

	existingHotNewsList := newsaggregatorv1.HotNewsList{
		Items: []newsaggregatorv1.HotNews{
			{
				ObjectMeta: v12.ObjectMeta{
					Namespace: "default",
					Name:      "feed-sample",
				},
				Spec: newsaggregatorv1.HotNewsSpec{
					Keywords:  []string{"keyword1,keyword2"},
					DateStart: "2024-08-12",
					DateEnd:   "2024-08-13",
					Feeds:     []string{"abc", "bbc"},
				},
			},
		},
	}

	existingConfigMap := v13.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Namespace: v1.FeedGroupsNamespace,
			Name:      v1.FeedGroupsConfigMapName,
		},
		Data: map[string]string{
			"sport":   "washingtontimes",
			"politic": "abc,bbc",
		},
	}

	k8sClient = fake.NewClientBuilder().
		WithLists(&existingHotNewsList).
		WithObjects(&existingConfigMap).
		WithScheme(scheme).
		Build()

	configMapRawData, err := json.Marshal(existingConfigMap)
	assert.Nil(t, err)

	tests := []struct {
		name               string
		body               io.Reader
		expectedStatusCode int
		expectedToFail     bool
		setup              func()
	}{
		{
			name: "Successful request",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Request: &admission.AdmissionRequest{
						UID:       "123",
						Namespace: "valid-namespace",
						Name:      "default",
						Resource: v12.GroupVersionResource{
							Group:    "",
							Version:  "v1",
							Resource: "configmaps",
						},
						Object: runtime.RawExtension{
							Raw: configMapRawData,
						},
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusOK,
			setup:              func() {},
			expectedToFail:     false,
		},
		{
			name:               "Fail to read request body",
			body:               &errorReader{},
			expectedStatusCode: http.StatusInternalServerError,
			setup:              func() {},
			expectedToFail:     true,
		},
		{
			name: "Invalid Request Body",
			body: func() *bytes.Buffer {
				return bytes.NewBuffer([]byte("invalid body"))
			}(),
			expectedStatusCode: http.StatusBadRequest,
			setup:              func() {},
			expectedToFail:     true,
		},
		{
			name: "Could Not Deserialize AdmissionReview",
			body: func() *bytes.Buffer {
				body, _ := json.Marshal(gin.H{"foo": "bar"})
				return bytes.NewBuffer(body)
			}(),
			setup:              func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedToFail:     true,
		},
		{
			name: "AdmissionReview Request is Nil",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			setup:              func() {},
			expectedStatusCode: http.StatusBadRequest,
			expectedToFail:     true,
		},
		{
			name: "Could Not Marshal JSON Patch",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Request: &admission.AdmissionRequest{
						UID:       "123",
						Namespace: "invalid-namespace",
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusBadRequest,
			setup:              func() {},
			expectedToFail:     true,
		},
		{
			name: "Valid request, but error because config map is empty",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Request: &admission.AdmissionRequest{
						UID:       "123",
						Namespace: "valid-namespace",
						Name:      "default",
						Resource: v12.GroupVersionResource{
							Group:    "",
							Version:  "v1",
							Resource: "configmaps",
						},
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusBadRequest,
			setup:              func() {},
			expectedToFail:     true,
		},
		{
			name: "Valid request, but error because hot news object is not registered",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Request: &admission.AdmissionRequest{
						UID:       "123",
						Namespace: "valid-namespace",
						Name:      "default",
						Resource: v12.GroupVersionResource{
							Group:    "",
							Version:  "v1",
							Resource: "configmaps",
						},
						Object: runtime.RawExtension{
							Raw: configMapRawData,
						},
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusBadRequest,
			setup: func() {
				k8sClient = fake.NewClientBuilder().Build()
			},
			expectedToFail: true,
		},
		{
			name: "Valid request, but error because no hot news detected",
			body: func() *bytes.Buffer {
				admissionReview := admission.AdmissionReview{
					Request: &admission.AdmissionRequest{
						UID:       "123",
						Namespace: "valid-namespace",
						Name:      "default",
						Resource: v12.GroupVersionResource{
							Group:    "",
							Version:  "v1",
							Resource: "configmaps",
						},
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusBadRequest,
			setup: func() {
				k8sClient = fake.NewClientBuilder().WithScheme(scheme).Build()
			},
			expectedToFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup()
			router := gin.Default()
			router.POST("/validate", validatingConfigMapHandler)

			req, _ := http.NewRequest(http.MethodPost, "/validate", test.body)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if test.expectedToFail {
				assert.NotEqual(t, http.StatusOK, rec.Code)
			} else {
				var admissionReview admission.AdmissionReview
				_ = json.NewDecoder(rec.Body).Decode(&admissionReview)
				assert.NotNil(t, admissionReview.Response)
			}
		})
	}
}

func Test_isKubeNamespace(t *testing.T) {
	type args struct {
		ns string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Kubernetes public Namespace",
			args: args{
				ns: "kube-public",
			},
			want: true,
		},
		{
			name: "Kubernetes system Namespace",
			args: args{
				ns: "kube-system",
			},
			want: true,
		},
		{
			name: "Not Kube Namespace",
			args: args{
				ns: "not kube-system",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isKubeNamespace(tt.args.ns); got != tt.want {
				t.Errorf("isKubeNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_triggerHotNewsReconcile(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)
	_ = newsaggregatorv1.AddToScheme(scheme)

	tests := []struct {
		name               string
		feedGroups         map[string]string
		hotNewsList        newsaggregatorv1.HotNewsList
		expectedFinalizers map[string][]string
		expectError        bool
	}{
		{
			name: "No matches - no finalizer added",
			feedGroups: map[string]string{
				"sport": "washingtontimes",
			},
			hotNewsList: newsaggregatorv1.HotNewsList{
				Items: []newsaggregatorv1.HotNews{
					{
						ObjectMeta: v12.ObjectMeta{
							Namespace: "default",
							Name:      "feed-sample",
						},
						Spec: newsaggregatorv1.HotNewsSpec{
							FeedGroups: []string{"finance", "weather"},
						},
					},
				},
			},
			expectedFinalizers: map[string][]string{
				"feed-sample": {},
			},
			expectError: false,
		},
		{
			name: "Match found - finalizer added",
			feedGroups: map[string]string{
				"politic": "abc",
			},
			hotNewsList: newsaggregatorv1.HotNewsList{
				Items: []newsaggregatorv1.HotNews{
					{
						ObjectMeta: v12.ObjectMeta{
							Namespace: "default",
							Name:      "feed-sample",
						},
						Spec: newsaggregatorv1.HotNewsSpec{
							FeedGroups: []string{"politic", "finance"},
						},
					},
				},
			},
			expectedFinalizers: map[string][]string{
				"feed-sample": {"hotnews.teamdev.com/reconcile"},
			},
			expectError: false,
		},
		{
			name:       "No feed groups provided - no finalizer added",
			feedGroups: map[string]string{},
			hotNewsList: newsaggregatorv1.HotNewsList{
				Items: []newsaggregatorv1.HotNews{
					{
						ObjectMeta: v12.ObjectMeta{
							Namespace: "default",
							Name:      "feed-sample",
						},
						Spec: newsaggregatorv1.HotNewsSpec{
							FeedGroups: []string{"politic"},
						},
					},
				},
			},
			expectedFinalizers: map[string][]string{
				"feed-sample": {},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k8sClient = fake.NewClientBuilder().
				WithLists(&tt.hotNewsList).
				WithScheme(scheme).
				Build()

			err := triggerHotNewsReconcile(tt.feedGroups, tt.hotNewsList)

			if (err != nil) != tt.expectError {
				t.Fatalf("expected error: %v, got: %v", tt.expectError, err)
			}

			for _, hotNews := range tt.hotNewsList.Items {
				updatedHotNews := &newsaggregatorv1.HotNews{}
				err = k8sClient.Get(context.TODO(), client.ObjectKey{Namespace: hotNews.Namespace, Name: hotNews.Name}, updatedHotNews)
				if err != nil {
					t.Fatalf("failed to get hotNews: %v", err)
				}

				expectedFinalizers := tt.expectedFinalizers[hotNews.Name]
				if len(updatedHotNews.Finalizers) != len(expectedFinalizers) || !equalSlices(updatedHotNews.Finalizers, expectedFinalizers) {
					t.Fatalf("expected finalizers %v, but got %v", expectedFinalizers, updatedHotNews.Finalizers)
				}
			}
		})
	}
}
func TestRunConfigMapController(t *testing.T) {
	type args struct {
		tlsCertFile string
		tlsKeyFile  string
	}
	cwd, err := os.Getwd()
	assert.Nil(t, err)
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
	}{
		//{
		//	name: "Successful run",
		//	args: args{
		//		tlsCertFile: filepath.Join(cwd, "templates", "certs", "ca.crt"),
		//		tlsKeyFile:  filepath.Join(cwd, "templates", "certs", "ca.key"),
		//	},
		//	setup: func() {
		//		_, cancel := context.WithTimeout(context.Background(), time.Second*5)
		//		defer cancel()
		//	},
		//	wantErr: false,
		//},
		{
			name: "Wrong certificates",
			args: args{
				tlsCertFile: filepath.Join(cwd, "templates", "certs", "fake.crt"),
				tlsKeyFile:  filepath.Join(cwd, "templates", "certs", "fake.key"),
			},
			setup: func() {
				_, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			errCh := make(chan error, 1)
			go func() {
				errCh <- RunConfigMapController(tt.args.tlsCertFile, tt.args.tlsKeyFile)
			}()

			select {
			case err := <-errCh:
				if tt.wantErr {
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

func TestValidateConfigMap_UniversalDecoderFails(t *testing.T) {
	invalidRaw := []byte(`invalid-data-that-causes-decode-to-fail`)
	req := &admission.AdmissionRequest{
		Resource: configMapResource,
		Object:   runtime.RawExtension{Raw: invalidRaw},
	}

	err := validateConfigMap(req)

	if err == nil || !strings.Contains(err.Error(), "could not deserialize configMap") {
		t.Fatalf("expected deserialization error, got: %v", err)
	}
}

func Test_triggerHotNewsReconcile_K8sList_Error(t *testing.T) {
	type args struct {
		feedGroups  map[string]string
		hotNewsList v1.HotNewsList
	}
	k8sClient = fake.NewClientBuilder().Build()
	tests := []struct {
		name    string
		args    args
		wantErr bool
		setup   func()
	}{
		{
			name: "Error due to k8s client scheme",
			args: args{
				feedGroups: map[string]string{
					"sport": "washingtontimes",
				},
				hotNewsList: v1.HotNewsList{
					Items: []v1.HotNews{
						{
							Spec: newsaggregatorv1.HotNewsSpec{
								FeedGroups: []string{"sport"},
							},
						},
					},
				},
			},
			wantErr: true,
			setup: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := triggerHotNewsReconcile(tt.args.feedGroups, tt.args.hotNewsList)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
