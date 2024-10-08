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
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"time"
)

var (
	feedlog = logf.Log.WithName("feed-resource")
)

// SetupWebhookWithManager will set up the manager to manage the webhooks
func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
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

	if r.OwnerReferences != nil {
		return nil, errors.New("feed has owner references")
	}

	return nil, nil
}

// validateFeed calls to our validation package to validate the feed configuration
func (r *Feed) validateFeed() (admission.Warnings, error) {
	var errList field.ErrorList

	err := validateFeedSpec(r.Spec)
	if err != nil {
		errList = append(errList, field.Invalid(field.NewPath("spec"), r.Spec, err.Error()))
	}

	err = r.checkNameUniqueness()
	if err != nil {
		errList = append(errList, field.Invalid(field.NewPath("spec.Name"), r.Spec, err.Error()))
	}

	err = r.checkLinkUniqueness()
	if err != nil {
		errList = append(errList, field.Invalid(field.NewPath("spec.Link"), r.Spec, err.Error()))
	}

	if len(errList) > 0 {
		return nil, errList.ToAggregate()
	}

	return nil, nil
}

// checkNameUniqueness checks if the Spec.name of the feed is unique in the namespace
func (r *Feed) checkNameUniqueness() error {
	feeds := &FeedList{}

	listOptions := client.ListOptions{Namespace: r.Namespace}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := k8sClient.List(ctx, feeds, &listOptions)
	if err != nil {
		return err
	}

	for _, feed := range feeds.Items {
		if feed.UID != r.UID {
			if feed.Spec.Name == r.Spec.Name && feed.Namespace == r.Namespace {
				return errors.New("name must be unique in the namespace")
			}
		}
	}

	return nil
}

// checkLinkUniqueness checks if the Spec.link of the feed is unique in the namespace
func (r *Feed) checkLinkUniqueness() error {
	feeds := &FeedList{}

	listOptions := client.ListOptions{Namespace: r.Namespace}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	err := k8sClient.List(ctx, feeds, &listOptions)
	if err != nil {
		return err
	}

	for _, feed := range feeds.Items {
		if feed.UID != r.UID {
			if feed.Spec.Link == r.Spec.Link && feed.Namespace == r.Namespace {
				return errors.New("link must be unique in the namespace")
			}
		}
	}

	return nil
}
