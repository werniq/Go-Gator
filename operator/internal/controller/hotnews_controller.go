/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

const (
	// feedGroupsNamespace is a namespace where feed groups are stored
	feedGroupsNamespace = "operator-system"

	// serverUrl is a URL to our news aggregator server
	serverUrl = "https://go-gator-svc.go-gator.svc.cluster.local:443/news"

	// feedGroupsConfigMapName is a name of the default ConfigMap which contains our feed groups names and sources
	feedGroupsConfigMapName = "feed-group-source"

	// errFeedsAreRequired is thrown when feeds are not provided
	errFeedsAreRequired = "feeds or feedGroups are required"

	// errFailedToConstructRequestUrl error message which is returned when failed to construct request URL
	errFailedToConstructRequestUrl = "failed to construct request URL"

	// errFailedToCreateRequest is returned when failed to create a new request
	errFailedToCreateRequest = "failed to create a new request"

	// errFailedToSendRequest indicates error during sending an HTTP request
	errFailedToSendRequest = "failed to send a request"

	// errFailedToDecodeResBody indicates that error occurred when failed to unmarshal response body
	errFailedToDecodeResBody = "failed to decode response body"

	// errFailedToCloseResponseBody is returned when failed to close response body
	errFailedToCloseResponseBody = "failed to close response body"

	// errWrongFeedGroupName is returned when the feed group name is wrong
	errWrongFeedGroupName = "wrong feed group name, please check the feed group name and try again"
)

// HotNewsReconciler reconciles a HotNews object
type HotNewsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews;configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/status;configmaps,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/finalizers;configmaps,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var hotNews newsaggregatorv1.HotNews

	logger.Info("Entering reconciler")

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		logger.Error(err, "unable to fetch HotNews")
		return ctrl.Result{}, err
	}

	if !hotNews.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("HotNews is deleting")
		err = r.handleDelete(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("HotNews has been deleted")
		return ctrl.Result{}, nil
	}

	if !hotNews.ObjectMeta.CreationTimestamp.IsZero() {
		logger.Info("HotNews is creating")
		err = r.handleCreate(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("HotNews has been created")
	} else {
		logger.Info("HotNews is updating")
		err = r.handleUpdate(ctx, &hotNews)
		if err != nil {
			return ctrl.Result{}, err
		}
		logger.Info("HotNews has been updated")
	}

	err = r.Client.Status().Update(ctx, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("HotNewsReconciler has been successfully executed")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}

// handleCreate function sends a request to the news aggregator server to retrieve news
// with the specified parameters, and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleCreate(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)

	requestUrl, err := r.constructRequestUrl(hotNews.Spec)
	if err != nil {
		logger.Error(err, errFailedToConstructRequestUrl)
		return err
	}
	logger.Info(requestUrl)

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		logger.Error(err, errFailedToCreateRequest)
		return err
	}

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	customClient := &http.Client{Transport: customTransport}

	res, err := customClient.Do(req)
	if err != nil {
		logger.Error(err, errFailedToSendRequest)
		return err
	}

	if res.StatusCode != http.StatusOK {
		serverErr := &serverError{}
		err = json.NewDecoder(res.Body).Decode(&serverErr)
		if err != nil {
			logger.Error(err, errFailedToDecodeResBody)
			return err
		}
		return serverErr
	}

	var articles struct {
		TotalNews int       `json:"totalAmount"`
		News      []article `json:"news"`
	}

	err = json.NewDecoder(res.Body).Decode(&articles)
	if err != nil {
		logger.Error(err, errFailedToDecodeResBody)
		return err
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(err, errFailedToCloseResponseBody)
		return err
	}

	var articlesTitles []string
	for _, a := range articles.News {
		articlesTitles = append(articlesTitles, a.Title)
	}
	logger.Info("Total amoutn of news", "totalAmount", articles.TotalNews)

	hotNews.InitHotNewsStatus(articles.TotalNews, requestUrl, articlesTitles)

	logger.Info("HotNews.handleCreate has been successfully executed")
	logger.Info("HotNews object", "HotNews", hotNews)

	return nil
}

// articles struct is used to parse the response from the news aggregator server
type article struct {
	Title       string `json:"title" xml:"title"`
	PubDate     string `json:"publishedAt" xml:"pubDate"`
	Description string `json:"description" xml:"description"`
	Publisher   string `xml:"source" json:"Publisher"`
	Link        string `json:"url" xml:"link"`
}

// handleUpdate function updates the HotNews object and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleUpdate(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)
	logger.Info("handling update")
	return nil
}

// handleDelete function deletes the HotNews object and returns an error if something goes wrong.
func (r *HotNewsReconciler) handleDelete(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)
	logger.Info("handling delete")
	return nil
}

// constructRequestUrl function verifies if arguments are correct and constructs a request URL
// to our news aggregator server.
//
// Example:
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&dateEnd=2024-08-06&sources=abc,bbc
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&sources=abc,bbc
func (r *HotNewsReconciler) constructRequestUrl(spec newsaggregatorv1.HotNewsSpec) (string, error) {
	var requestUrl strings.Builder

	requestUrl.WriteString(serverUrl)
	requestUrl.WriteString("?keywords=" + spec.Keywords)

	if len(spec.Feeds) < 1 && spec.FeedGroups == nil {
		return "", fmt.Errorf(errFeedsAreRequired)
	}

	var feedStr strings.Builder
	if spec.FeedGroups != nil {
		feedGroups, err := r.processFeedGroups(spec)
		if err != nil {
			return "", err
		}
		for _, feedGroup := range strings.Split(feedGroups, ",") {
			feedStr.WriteString(feedGroup)
			feedStr.WriteRune(',')
		}
	} else if spec.Feeds != nil {
		for _, feed := range spec.Feeds {
			feedStr.WriteString(feed)
			feedStr.WriteRune(',')
		}
	}
	if feedStr.String() == "" {
		requestUrl.WriteString("&sources=" + feedStr.String()[:len(feedStr.String())-1])
	}

	if spec.DateStart != "" {
		requestUrl.WriteString("&dateFrom=" + spec.DateStart)
	}

	if spec.DateEnd != "" {
		requestUrl.WriteString("&dateEnd=" + spec.DateEnd)
	}

	return requestUrl.String(), nil
}

// processFeedGroups function processes feed groups from the ConfigMap and returns a string with feed sources
func (r *HotNewsReconciler) processFeedGroups(spec newsaggregatorv1.HotNewsSpec) (string, error) {
	var sources strings.Builder

	feedGroups, err := r.getFeedGroups(context.Background())
	if err != nil {
		return "", err
	}

	for _, feedKey := range spec.FeedGroups {
		if source, exists := feedGroups.Data[feedKey]; exists {
			sources.WriteString(source)
			sources.WriteRune(',')
		} else {
			return "", fmt.Errorf(errWrongFeedGroupName)
		}
	}

	return sources.String(), nil
}

// getConfigMapData returns all data from config map named feedGroupsConfigMapName in defaultNamespace
func (r *HotNewsReconciler) getFeedGroups(ctx context.Context) (*v1.ConfigMap, error) {
	c := config.GetConfigOrDie()

	k8sClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	configMap, err := k8sClient.CoreV1().ConfigMaps(feedGroupsNamespace).
		Get(ctx, feedGroupsConfigMapName, v12.GetOptions{})
	if err != nil {
		return nil, err
	}

	return configMap, nil
}
