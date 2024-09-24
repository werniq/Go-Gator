package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	admission "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"log"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	"time"
)

var (
	// universalDeserializer is a deserializer for Kubernetes objects, used to decode the incoming HTTP request
	universalDeserializer = runtime.Decoder(nil)

	// configMapResource represents the Kubernetes resource type for a config map
	configMapResource = metav1.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}

	// k8sClient is the Kubernetes client used to interact with the API server, used to retrieve all hotnews from the
	// config map's namespace, and to trigger a reconcile of all hotnews which have the feed group in their feed groups.
	k8sClient client.Client

	scheme = runtime.NewScheme()
)

const (
	// errEmptyData identifies an error when the configMap data is nil, and no feed groups to reconcile
	errEmptyData = "configMap data is nil, so no feed groups to reconcile"
)

// RunConfigMapController starts the admission controller webhook for validating config maps.
//
// It listens on port 8443 and delegates the admission control logic to the validatingConfigMapHandler func.
func RunConfigMapController(tlsCertFile, tlsKeyFile string) error {
	var err error

	k8sClient, err = client.New(ctrl.GetConfigOrDie(), client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return err
	}

	universalDeserializer = serializer.NewCodecFactory(scheme).UniversalDeserializer()

	r := gin.Default()

	setupRoutes(r)

	err = r.RunTLS(":8443", tlsCertFile, tlsKeyFile)
	if err != nil {
		return err
	}

	return nil
}

// webhookApiResponse is minimal required response from a webhook to allow or forbid a request
type webhookApiResponse struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Response   response `json:"response"`
}

// response struct contains Uid and Allowed fields, which describe if webhook has validated succesffully, or not.
type response struct {
	// Uid is used to match the response to the corresponding request
	Uid types.UID `json:"uid"`

	// Allowed field is a boolean value, which is true if the request is allowed, and false otherwise.
	// If webhook returned false, the request will be rejected.
	Allowed bool `json:"allowed"`

	// Status field contains the HTTP status code and a message, which is returned in case of an error.
	Status *status `json:"status,omitempty"`
}

// status struct contains the HTTP status code and a message, which is returned in case of an error.
// It will make the response more informative and user-friendly.
type status struct {
	// Code field contains the HTTP status code, which is either 200 (OK) or 400 (Bad Request)
	Code int `json:"code"`

	// Message field contains a message, which is used to describe the error in case of a bad request,
	// or to provide additional information in case of a successful request.
	Message string `json:"message"`
}

// setupRoutes configures the HTTP routes for the admission controller webhook.
func setupRoutes(r *gin.Engine) {
	r.POST("/validate", validatingConfigMapHandler)
}

// validatingConfigMapHandler parses the HTTP request for an admission controller webhook, and -- in case of a well-formed
// request -- delegates the admission control logic to the validateConfigMap func
func validatingConfigMapHandler(c *gin.Context) {
	var err error

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		res := webhookApiResponse{
			ApiVersion: "admission.k8s.io/v1",
			Kind:       "configmaps",
			Response: response{
				Allowed: false,
				Status: &status{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("could not read request body: %v", err),
				},
			},
		}
		c.JSON(http.StatusInternalServerError, res)
		log.Printf("could not read request body: %v\n", err)
		return
	}

	var admissionReviewReq admission.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil,
		&admissionReviewReq); err != nil {
		res := webhookApiResponse{
			ApiVersion: "admission.k8s.io/v1",
			Kind:       "configmaps",
			Response: response{
				Allowed: false,
				Status: &status{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("could not deserialize request: %v", err),
				},
			},
		}
		c.JSON(http.StatusBadRequest, res)
		return
	} else if admissionReviewReq.Request == nil {
		res := webhookApiResponse{
			ApiVersion: "admission.k8s.io/v1",
			Kind:       "configmaps",
			Response: response{
				Allowed: false,
				Status: &status{
					Code:    http.StatusBadRequest,
					Message: "malformed admission review: request is nil",
				},
			},
		}
		c.JSON(http.StatusBadRequest, res)
		log.Println("malformed admission review: request is nil")
		return
	}

	res := webhookApiResponse{
		ApiVersion: admissionReviewReq.APIVersion,
		Kind:       admissionReviewReq.Kind,
		Response: response{
			Uid: admissionReviewReq.Request.UID,
		},
	}

	if !isKubeNamespace(admissionReviewReq.Request.Namespace) {
		err = validateConfigMap(admissionReviewReq.Request)
	}

	if err != nil {
		res.Response.Allowed = false
		res.Response.Status = &status{
			Message: fmt.Sprintf("Error during validation of configMap: %v", err),
			Code:    http.StatusBadRequest,
		}
		log.Println("Error during validation of configMap: ", err)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res = webhookApiResponse{
		ApiVersion: admissionReviewReq.APIVersion,
		Kind:       admissionReviewReq.Kind,
		Response: response{
			Uid:     admissionReviewReq.Request.UID,
			Allowed: true,
			Status: &status{
				Code: http.StatusOK,
			},
		},
	}

	c.JSON(http.StatusOK, res)
}

// validateConfigMap verifies that the configMap has a data field and triggers a reconcile of all hotnews
// which have the feed group in their feed groups.
func validateConfigMap(req *admission.AdmissionRequest) error {
	if req.Resource != configMapResource {
		return fmt.Errorf("expect resource to be %s, got %s", configMapResource, req.Resource)
	}

	raw := req.Object.Raw
	configMap := v1.ConfigMap{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &configMap); err != nil {
		return fmt.Errorf("could not deserialize configMap: %v", err)
	}

	if configMap.Data == nil {
		return fmt.Errorf(errEmptyData)
	}

	feeds, err := getAllHotNewsFromNamespace(configMap.Namespace)
	if err != nil {
		return err
	}

	err = triggerHotNewsReconcile(configMap.Data, feeds)
	if err != nil {
		return err
	}

	return nil
}

// getAllHotNewsFromNamespace retrieves all hotnews from the provided namespace
func getAllHotNewsFromNamespace(namespace string) (newsaggregatorv1.HotNewsList, error) {
	var hotNewsList newsaggregatorv1.HotNewsList
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := k8sClient.List(ctx, &hotNewsList, client.InNamespace(namespace))
	if err != nil {
		return newsaggregatorv1.HotNewsList{}, fmt.Errorf("error retrieving hotnews CRD: %v\n", err)
	}

	return hotNewsList, nil
}

// triggerHotNewsReconcile triggers a reconcile of all hotnews which have the given feed group in their feed groups.
func triggerHotNewsReconcile(feedGroups map[string]string, hotNewsList newsaggregatorv1.HotNewsList) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for feedGroup, _ := range feedGroups {
		for _, hotNews := range hotNewsList.Items {
			if slices.Contains(hotNews.Spec.FeedGroups, feedGroup) {
				hotNews.Finalizers = append(hotNews.Finalizers, "hotnews.teamdev.com/reconcile")
				err := k8sClient.Update(ctx, &hotNews)
				if err != nil {
					return fmt.Errorf("could not update hotNews %s: %v\n", hotNews.Name, err)
				}
			}
		}
	}

	return nil
}

// isKubeNamespace checks if the given namespace is a Kubernetes-owned namespace.
func isKubeNamespace(ns string) bool {
	return ns == metav1.NamespacePublic || ns == metav1.NamespaceSystem
}
