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

package v1

import (
	"context"
	"fmt"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// errNoFeeds is an error message indicating that user hasn't specified any feeds
	errNoFeeds = "either feeds or feedGroups should be specified"

	// errInvalidDateRange is an error message indicating input of wrong date range
	errInvalidDateRange = "DateStart should be before than DateEnd"

	// errWrongFeedGroupName is an error message for wrong hotNew group name
	errWrongFeedGroupName = "hotNew group name is not found in the ConfigMap, please check the hotNew group name"

	// FeedGroupsNamespace is a namespace where hotNew groups are stored
	FeedGroupsNamespace = "operator-system"

	// FeedGroupsConfigMapName is a name of the default ConfigMap which contains our hotNew groups names and sources
	FeedGroupsConfigMapName = "feed-group-source"
)

var (
	hotnewslog = logf.Log.WithName("hotnews-resource")

	// c is a kubernetes configuration which will be used to create a k8s client
	c = config.GetConfigOrDie()

	// k8sClient is a k8s client which will be used to get ConfigMap with hotNew groups
	clientset *kubernetes.Clientset
)

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *HotNews) SetupWebhookWithManager(mgr ctrl.Manager) error {
	var err error
	clientset, err = kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-newsaggregator-teamdev-com-v1-hotnews,mutating=true,failurePolicy=fail,sideEffects=None,groups=newsaggregator.teamdev.com,resources=hotnews,verbs=create;update;delete,versions=v1,name=mhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &HotNews{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
//
// This webhook will set the default values for the HotNews resource
// In particular, if the user hasn't specified the number of titles to show in the summary, we will set it to 10
func (r *HotNews) Default() {
	hotnewslog.Info("default", "name", r.Name)

	var err error
	var feeds []Feed

	if r.Spec.SummaryConfig.TitlesCount == 0 {
		r.Spec.SummaryConfig.TitlesCount = 10
	}

	if r.Spec.Feeds == nil && r.Spec.FeedGroups == nil {
		feeds, err = r.getAllFeedsInCurrentNamespace(context.TODO())
	}

	for _, feed := range feeds {
		r.Spec.Feeds = append(r.Spec.Feeds, feed.Spec.Name)
	}

	if err != nil {
		log.Fatalln(err)
	}
}

// +kubebuilder:webhook:path=/validate-newsaggregator-teamdev-com-v1-hotnews,mutating=false,failurePolicy=fail,sideEffects=None,groups=newsaggregator.teamdev.com,resources=hotnews,verbs=create;update;delete,versions=v1,name=vhotnews.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &HotNews{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
//
// It is called when the HotNews resource is created
// Validating webhook will check if the HotNews resource is correct
// In particular, it checks if the DateStart is before DateEnd and if all hotNew group names are correct
// Also, it checks if user-specified feeds or feedGroups are correct by these criteria:
// FeedGroups should be present in the feed-group-source ConfigMap
func (r *HotNews) ValidateCreate() (admission.Warnings, error) {
	hotnewslog.Info("validate create", "name", r.Name)
	err := r.validateHotNews()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
//
// ValidateUpdate is called when the HotNews resource is Updated
// Validating webhook will check if the HotNews resource is correct
// In particular, it checks if the DateStart is before DateEnd and if all hotNew group names are correct
// Also, it checks if user-specified feeds or feedGroups are correct by these criteria:
// FeedGroups should be present in the feed-group-source ConfigMap
func (r *HotNews) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	hotnewslog.Info("validate update", "name", r.Name)
	err := r.validateHotNews()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *HotNews) ValidateDelete() (admission.Warnings, error) {
	hotnewslog.Info("validate delete", "name", r.Name)

	return nil, nil
}

// validateHotNews validates the HotNews resource.
//
// In particular, it checks if the DateStart is before DateEnd and if all hotNew group names are correct, and
// if feeds or feedGroups exists in our news aggregator.
func (r *HotNews) validateHotNews() error {
	if r.Spec.DateStart > r.Spec.DateEnd {
		return fmt.Errorf(errInvalidDateRange)
	}

	if r.Spec.Feeds == nil && r.Spec.FeedGroups == nil {
		return fmt.Errorf(errNoFeeds)
	}

	configMap, err := clientset.CoreV1().ConfigMaps(FeedGroupsNamespace).
		Get(context.TODO(), FeedGroupsConfigMapName, v12.GetOptions{})

	if err != nil {
		hotnewslog.Info("Error retireving config map")
		return err
	}

	for _, source := range r.Spec.FeedGroups {
		if _, exists := configMap.Data[source]; !exists {
			return fmt.Errorf(errWrongFeedGroupName)
		}
	}

	return nil
}

// getAllFeedsInCurrentNamespace returns all feeds in the current namespace
func (r *HotNews) getAllFeedsInCurrentNamespace(ctx context.Context) ([]Feed, error) {
	var feeds FeedList
	err := k8sClient.List(ctx, &feeds)
	if err != nil {
		return nil, err
	}

	return feeds.Items, nil
}
