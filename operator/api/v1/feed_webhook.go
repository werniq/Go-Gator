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
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	feedlog = logf.Log.WithName("feed-resource")

	// k8sClient is a kubernetes client that is used to interact with the k8s API
	k8sClient client.Client
)

// SetupWebhookWithManager will set up the manager to manage the webhooks
func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	var err error
	clientset, err = kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}

	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-newsaggregator-teamdev-com-v1-feed,mutating=false,failurePolicy=fail,sideEffects=None,groups=newsaggregator.teamdev.com,resources=feeds,verbs=create;update;delete,versions=v1,name=vfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Feed{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateCreate() (admission.Warnings, error) {
	feedlog.Info("validate create", "name", r.Name)
	return r.validateFeed()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	feedlog.Info("validate update", "name", r.Name)
	return r.validateFeed()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Feed) ValidateDelete() (admission.Warnings, error) {
	feedlog.Info("validate delete", "name", r.Name)
	return nil, nil
}

// validateFeed calls to our validation package to validate the feed configuration
func (r *Feed) validateFeed() (admission.Warnings, error) {
	err := validateFeeds(r.Spec)
	if err != nil {
		return nil, err
	}

	warn, err := r.checkNameUniqueness()
	if err != nil {
		return warn, err
	}

	warn, err = r.checkLinkUniqueness()
	if err != nil {
		return warn, err
	}

	return nil, nil
}

// checkNameUniqueness checks if the Spec.name of the feed is unique in the namespace
func (r *Feed) checkNameUniqueness() (admission.Warnings, error) {
	feeds := &FeedList{}

	listOptions := client.ListOptions{Namespace: r.Namespace}

	err := k8sClient.List(context.Background(), feeds, &listOptions)
	if err != nil {
		return nil, err
	}

	for _, feed := range feeds.Items {
		if feed.Spec.Name == r.Spec.Name && feed.Namespace == r.Namespace {
			return nil, errors.New("name must be unique in the namespace")
		}
	}

	return nil, nil
}

// checkLinkUniqueness checks if the Spec.link of the feed is unique in the namespace
func (r *Feed) checkLinkUniqueness() (admission.Warnings, error) {
	feeds := &FeedList{}

	listOptions := client.ListOptions{Namespace: r.Namespace}

	err := k8sClient.List(context.Background(), feeds, &listOptions)
	if err != nil {
		return nil, err
	}

	for _, feed := range feeds.Items {
		if feed.Spec.Link == r.Spec.Link && feed.Namespace == r.Namespace {
			return nil, errors.New("link must be unique in the namespace")
		}
	}

	return nil, nil
}
