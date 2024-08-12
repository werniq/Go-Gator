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
	"encoding/json"
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	config "sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"teamdev.com/go-gator/api/v1/validation"
)

// log is for logging in this package.
var feedlog = logf.Log.WithName("feed-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *Feed) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-newsaggregator-teamdev-com-v1-feed,mutating=true,failurePolicy=fail,sideEffects=None,groups=newsaggregator.teamdev.com,resources=feeds,verbs=create;update;delete,versions=v1,name=mfeed.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Feed{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Feed) Default() {
	feedlog.Info("default", "name", r.Name)

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
	err := validation.Validate(r.Spec.Name, r.Spec.Link)
	if err != nil {
		return nil, err
	}

	c := config.GetConfigOrDie()
	k8sClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	feeds := FeedList{}
	d, err := k8sClient.RESTClient().
		Get().
		AbsPath("/apis/newsaggregator.teamdev.com/v1/feeds").
		DoRaw(context.TODO())

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(d, &feeds)
	if err != nil {
		return nil, err
	}

	for _, feed := range feeds.Items {
		if feed.Spec.Name == r.Spec.Name {
			return nil, errors.New("name must be unique")
		}
	}

	return nil, nil
}
