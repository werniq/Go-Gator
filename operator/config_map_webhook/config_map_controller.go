package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	admission "k8s.io/api/admission/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
	time "time"
)

var (
	// universalDeserializer is a deserializer for Kubernetes objects, used to decode the incoming HTTP request
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()

	// configMapResource represents the Kubernetes resource type for a config map
	configMapResource = metav1.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}

	// k8sClient is the Kubernetes client used to interact with the API server, used to retrieve all hotnews from the
	// config map's namespace, and to trigger a reconcile of all hotnews which have the feed group in their feed groups.
	k8sClient client.Client
)

// patchOperation is an operation of a JSON patch, see https://tools.ietf.org/html/rfc6902 .
type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// RunConfigMapController starts the admission controller webhook for validating config maps.
//
// It listens on port 8732 and delegates the admission control logic to the validatingConfigMapHandler func.
func RunConfigMapController(client client.Client) error {
	k8sClient = client

	r := gin.Default()

	setupRoutes(r)

	err := r.RunTLS(":8443", tlsCertFile, tlsKeyFile)
	if err != nil {
		return err
	}

	return nil
}

// setupRoutes configures the HTTP routes for the admission controller webhook.
func setupRoutes(r *gin.Engine) {
	r.POST("/validate", validatingConfigMapHandler)
}

// validatingConfigMapHandler parses the HTTP request for an admission controller webhook, and -- in case of a well-formed
// request -- delegates the admission control logic to the validateConfigMap func
//
// The response body is then returned as raw bytes.
func validatingConfigMapHandler(c *gin.Context) {
	var err error

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not read request body: %v", err)})
		log.Printf("could not read request body: %v\n", err)
		return
	}

	var admissionReviewReq admission.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &admissionReviewReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not deserialize request: %v", err)})
		log.Printf("could not deserialize request: %v\n", err)
		return
	} else if admissionReviewReq.Request == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "malformed admission review: request is nil"})
		log.Println("malformed admission review: request is nil")
		return
	}

	admissionReviewResponse := admission.AdmissionReview{
		TypeMeta: admissionReviewReq.TypeMeta,
		Response: &admission.AdmissionResponse{
			UID: admissionReviewReq.Request.UID,
		},
	}

	var patchOps []patchOperation
	if !isKubeNamespace(admissionReviewReq.Request.Namespace) {
		patchOps, err = validateConfigMap(admissionReviewReq.Request)
	}

	if err != nil {
		admissionReviewResponse.Response.Allowed = false
		admissionReviewResponse.Response.Result = &metav1.Status{
			Message: err.Error(),
		}
	} else {
		patchBytes, err := json.Marshal(patchOps)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not marshal JSON patch: %v", err)})
			log.Printf("could not marshal JSON patch: %v\n", err)
			return
		}
		admissionReviewResponse.Response.Allowed = true
		admissionReviewResponse.Response.Patch = patchBytes

		admissionReviewResponse.Response.PatchType = new(admission.PatchType)
		*admissionReviewResponse.Response.PatchType = admission.PatchTypeJSONPatch
	}

	bytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		log.Printf("marshaling response: %v\n", err)
	}

	c.JSON(http.StatusOK, bytes)
}

// validateConfigMap verifies that the configMap has a data field and triggers a reconcile of all hotnews
// which have the feed group in their feed groups.
func validateConfigMap(req *admission.AdmissionRequest) ([]patchOperation, error) {
	if req.Resource != configMapResource {
		return nil, fmt.Errorf("expect resource to be %s, got %s", configMapResource, req.Resource)
	}

	raw := req.Object.Raw
	configMap := v1.ConfigMap{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &configMap); err != nil {
		return nil, fmt.Errorf("could not deserialize configMap: %v", err)
	}

	if configMap.Data == nil {
		return nil, fmt.Errorf("configMap data is nil, so no feed groups to reconcile")
	}

	feeds, err := getAllHotNewsFromNamespace(configMap.Namespace)
	if err != nil {
		return nil, err
	}

	triggerHotNewsReconcile(configMap.Data, feeds)

	return nil, nil
}

// getAllHotNewsFromNamespace retrieves all hotnews from the provided namespace
func getAllHotNewsFromNamespace(namespace string) (newsaggregatorv1.HotNewsList, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// retrieve all hotnews from the config map's namespace
	var feeds newsaggregatorv1.HotNewsList
	err := k8sClient.List(ctx, &feeds, client.InNamespace(namespace))
	if err != nil {
		return newsaggregatorv1.HotNewsList{}, fmt.Errorf("could not get feeds: %v", err)
	}

	return feeds, err
}

// triggerHotNewsReconcile triggers a reconcile of all hotnews which have the given feed group in their feed groups.
func triggerHotNewsReconcile(feedGroups map[string]string, feeds newsaggregatorv1.HotNewsList) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, feedGroup := range feedGroups {
		for _, feed := range feeds.Items {
			if slices.Contains(feed.Spec.FeedGroups, feedGroup) {
				feed.Finalizers = append(feed.Finalizers, "hotnews.teamdev.com/reconcile")
				err := k8sClient.Update(ctx, &feed)
				if err != nil {
					log.Printf("could not update feed %s: %v\n", feed.Name, err)
				}
			}
		}
	}
}

// isKubeNamespace checks if the given namespace is a Kubernetes-owned namespace.
func isKubeNamespace(ns string) bool {
	return ns == metav1.NamespacePublic || ns == metav1.NamespaceSystem
}
