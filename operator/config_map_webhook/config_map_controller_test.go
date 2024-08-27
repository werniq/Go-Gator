package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	admission "k8s.io/api/admission/v1beta1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"net/http"
	"net/http/httptest"
	"reflect"
	v1 "teamdev.com/go-gator/api/v1"
	"testing"
)

func TestWebhookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	universalDeserializer = serializer.NewCodecFactory(scheme).UniversalDeserializer()

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

	configMapRaw, err := json.Marshal(existingConfigMap)
	assert.Nil(t, err)

	tests := []struct {
		name               string
		body               *bytes.Buffer
		expectedStatusCode int
		expectedToFail     bool
	}{
		{
			name: "Invalid Request Body",
			body: func() *bytes.Buffer {
				return bytes.NewBuffer([]byte("invalid body"))
			}(),
			expectedStatusCode: http.StatusBadRequest,
			expectedToFail:     true,
		},
		{
			name: "Could Not Deserialize AdmissionReview",
			body: func() *bytes.Buffer {
				body, _ := json.Marshal(gin.H{"foo": "bar"})
				return bytes.NewBuffer(body)
			}(),
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
			expectedToFail:     true,
		},
		{
			name: "Valid request, but error retrieving hot news from namespace",
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
							Raw: configMapRaw,
						},
					},
					Response: &admission.AdmissionResponse{},
				}
				body, _ := json.Marshal(admissionReview)
				return bytes.NewBuffer(body)
			}(),
			expectedStatusCode: http.StatusOK,
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
			expectedStatusCode: http.StatusOK,
			expectedToFail:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			router.POST("/validate", validatingConfigMapHandler)

			req, _ := http.NewRequest(http.MethodPost, "/validate", test.body)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, test.expectedStatusCode, rec.Code)

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

func Test_getAllHotNewsFromNamespace(t *testing.T) {
	type args struct {
		namespace string
	}
	tests := []struct {
		name    string
		args    args
		want    v1.HotNewsList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAllHotNewsFromNamespace(tt.args.namespace)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAllHotNewsFromNamespace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAllHotNewsFromNamespace() got = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isKubeNamespace(tt.args.ns); got != tt.want {
				t.Errorf("isKubeNamespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setupRoutes(t *testing.T) {
	type args struct {
		r *gin.Engine
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRoutes(tt.args.r)
		})
	}
}

func Test_triggerHotNewsReconcile(t *testing.T) {
	type args struct {
		feedGroups  map[string]string
		hotNewsList v1.HotNewsList
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := triggerHotNewsReconcile(tt.args.feedGroups, tt.args.hotNewsList); (err != nil) != tt.wantErr {
				t.Errorf("triggerHotNewsReconcile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateConfigMap(t *testing.T) {
	type args struct {
		req *admission.AdmissionRequest
	}
	tests := []struct {
		name    string
		args    args
		want    []patchOperation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateConfigMap(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfigMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateConfigMap() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validatingConfigMapHandler(t *testing.T) {
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validatingConfigMapHandler(tt.args.c)
		})
	}
}
