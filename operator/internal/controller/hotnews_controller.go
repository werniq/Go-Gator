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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"net/http"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
	"time"

	newsaggregatorv1 "teamdev.com/go-gator/api/v1"
)

var (
	// c is a kubernetes configuration which will be used to create a k8s client
	c = config.GetConfigOrDie()

	// k8sClient is a k8s client which will be used to get ConfigMap with hotNew groups
	clientset *kubernetes.Clientset
)

const (
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
// Whenever status of HotNews CRD is updated, it sends a request to the news aggregator server
// to retrieve news with the specified parameters.
// It also watches for changes in the ConfigMap with the feed groups, and in the Feed CRD.
//
// Before sending a request to the news aggregator server, it verifies if the arguments are correct:
// - keywords are provided
// - date range is correct
// - feeds or feed groups are provided, and they exists in news aggregator server
// Then, it constructs a request URL and sends a request to the news aggregator server, parses the response
// and updates the HotNews object.
type HotNewsReconciler struct {
	serverUrl string
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews;feeds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/status;feeds,verbs=get;update;patch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=hotnews/finalizers;feeds,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:rbac:groups=newsaggregator.teamdev.com,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state
//
// This function will be called when a HotNews object is created, updated or deleted
// It will send a request to the news aggregator server to retrieve news with the parameters,
// specified in the HotNews object.
// Additionally, it is watching for changes in the ConfigMap with the feed groups, and in the Feed CRD.
// If there were any changes, it will also affect the HotNews object.
func (r *HotNewsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var hotNews newsaggregatorv1.HotNews

	err := r.Client.Get(ctx, req.NamespacedName, &hotNews)
	if err != nil {
		logger.Error(err, "unable to fetch HotNews")
		return ctrl.Result{}, err
	}

	err = r.processHotNews(ctx, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}
	logger.Info("HotNews object has been updated")

	err = r.Client.Status().Update(ctx, &hotNews)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger.Info("HotNewsReconciler has been successfully executed")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager, and initializes the k8s client
// to work with feedGroup Config Map.
// It Watches for any changes in the ConfigMap with the feed groups, and also watches for changes in Feed CRD.
func (r *HotNewsReconciler) SetupWithManager(mgr ctrl.Manager, serverUrl string) error {
	var err error
	clientset, err = kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}

	r.serverUrl = serverUrl

	return ctrl.NewControllerManagedBy(mgr).
		For(&newsaggregatorv1.HotNews{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(&v1.ConfigMap{
			ObjectMeta: v12.ObjectMeta{
				Name:      newsaggregatorv1.FeedGroupsConfigMapName,
				Namespace: newsaggregatorv1.FeedGroupsNamespace,
			},
		},
			&handler.EnqueueRequestForObject{}).
		Watches(&newsaggregatorv1.Feed{},
			&handler.EnqueueRequestForObject{}).
		Complete(r)
}

// articles struct is used to parse the response from the news aggregator server
type article struct {
	Title       string `json:"title" xml:"title"`
	PubDate     string `json:"publishedAt" xml:"pubDate"`
	Description string `json:"description" xml:"description"`
	Publisher   string `xml:"source" json:"Publisher"`
	Link        string `json:"url" xml:"link"`
}

// processHotNews function updates the HotNews object and returns an error if something goes wrong.
func (r *HotNewsReconciler) processHotNews(ctx context.Context, hotNews *newsaggregatorv1.HotNews) error {
	logger := log.FromContext(ctx)
	logger.Info("handling update")

	requestUrl, err := r.constructRequestUrl(ctx, hotNews.Spec)
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
		serverError := &serverErr{}
		err = json.NewDecoder(res.Body).Decode(&serverError)
		if err != nil {
			logger.Error(err, errFailedToDecodeResBody)
			return err
		}
		return serverError
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
	logger.Info("Total amount of news", "totalAmount", articles.TotalNews)

	hotNews.InitHotNewsStatus(articles.TotalNews, requestUrl, articlesTitles)

	logger.Info("HotNews.processHotNews has been successfully executed")
	logger.Info("HotNews object", "HotNews", hotNews)

	return nil
}

// constructRequestUrl function verifies if arguments are correct and constructs a request URL
// to our news aggregator server.
//
// Example:
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&dateEnd=2024-08-06&sources=abc,bbc
// http://server.com/news?keywords=bitcoin&dateFrom=2024-08-05&sources=abc,bbc
func (r *HotNewsReconciler) constructRequestUrl(ctx context.Context, spec newsaggregatorv1.HotNewsSpec) (string, error) {
	var requestUrl strings.Builder

	requestUrl.WriteString(r.serverUrl)
	var keywordsStr strings.Builder
	for _, keyword := range spec.Keywords {
		keywordsStr.WriteString(keyword)
		keywordsStr.WriteRune(',')
	}
	requestUrl.WriteString("?keywords=" + keywordsStr.String()[:len(keywordsStr.String())-1])

	var feedStr strings.Builder
	if spec.FeedGroups != nil {
		feedGroupsStr, err := r.processFeedGroups(spec)
		if err != nil {
			return "", err
		}
		feedStr.WriteString(feedGroupsStr)
	} else {
		feedsStr, err := r.processFeeds(spec)
		if err != nil {
			return "", err
		}
		feedStr.WriteString(feedsStr)
	}

	requestUrl.WriteString("&sources=" + feedStr.String())

	if spec.DateStart != "" {
		requestUrl.WriteString("&dateFrom=" + spec.DateStart)
	}

	if spec.DateEnd != "" {
		requestUrl.WriteString("&dateEnd=" + spec.DateEnd)
	}

	return requestUrl.String(), nil
}

// processFeeds returns a string containing comma-separated feed sources
func (r *HotNewsReconciler) processFeeds(spec newsaggregatorv1.HotNewsSpec) (string, error) {
	var sourcesBuilder strings.Builder

	for _, feed := range spec.Feeds {
		sourcesBuilder.WriteString(feed)
		sourcesBuilder.WriteRune(',')
	}

	return sourcesBuilder.String()[:len(sourcesBuilder.String())-1], nil
}

// processFeedGroups function processes feed groups from the ConfigMap and returns a string containing comma-separated
// feed sources
func (r *HotNewsReconciler) processFeedGroups(spec newsaggregatorv1.HotNewsSpec) (string, error) {
	var sourcesBuilder strings.Builder

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	configMaps, err := r.getFeedGroups(ctx)
	if err != nil {
		return "", err
	}

	for _, feedGroup := range spec.FeedGroups {
		for _, configMap := range configMaps.Items {
			if _, ok := configMap.Data[feedGroup]; !ok {
				return "", fmt.Errorf(errWrongFeedGroupName)
			} else {
				sourcesBuilder.WriteString(configMap.Data[feedGroup])
				sourcesBuilder.WriteRune(',')
			}
		}
	}

	return sourcesBuilder.String()[:len(sourcesBuilder.String())-1], nil
}

// getConfigMapData returns all data from config map named FeedGroupsConfigMapName in FeedGroupsNamespace
func (r *HotNewsReconciler) getFeedGroups(ctx context.Context) (v1.ConfigMapList, error) {
	var configMaps v1.ConfigMapList
	err := r.Client.List(ctx, &configMaps, client.InNamespace(newsaggregatorv1.FeedGroupsNamespace))

	if err != nil {
		return v1.ConfigMapList{}, err
	}

	logger := log.FromContext(ctx)
	for _, configMap := range configMaps.Items {
		for key, item := range configMap.Data {
			logger.Info("ConfigMap data", key, item)
		}
	}

	return configMaps, nil
}

// getAllFeedsInCurrentNamespace returns all feeds in the current namespace
func (r *HotNewsReconciler) getAllFeedsInCurrentNamespace(ctx context.Context) ([]newsaggregatorv1.Feed, error) {
	var feeds newsaggregatorv1.FeedList
	err := r.Client.List(ctx, &feeds)
	if err != nil {
		return nil, err
	}

	return feeds.Items, nil
}

// feedInNamespace returns true if feed is in the namespace, otherwise - false
func (r *HotNewsReconciler) feedInNamespace(namespace []string, feed string) bool {
	for _, source := range namespace {
		if source == feed {
			return true
		}
	}
	return false
}
