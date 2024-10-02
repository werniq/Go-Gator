/*
Package v1 contains CRD definitions for the Feed and HotNews resources. Validation and mutating webhooks that ensures uniqueness of the
feed.

This package includes the following components:
- Custom Resource Definition for the Feed resource.
- Webhook for validating Feed resources.
- CRD for the HotNews resource.
- Validating and mutating webhooks for the HotNews resource.

The Feed resource represents a news feed with a name and a link.

The webhook ensures that the Feed resources meet the required validation criteria, in particular:
- Name should be unique within namespace, and should not be more than 20 symbols
- Endpoint of feed should be either http or https url.

The HotNews resource represents a group of news feeds with a name and a list of feed names.
This CRD allows us to create hot news based on available Feeds.

The mutating webhook sets the default values for the HotNews resource, if they are not specified, like:
- Spec.SummaryConfig.TitlesCount should be 10 by default.

The validating webhook ensures that the HotNews resources meet the required validation criteria, in particular:
- Either Feeds or FeedGroups should be specified.
- DateStart should be before DateEnd.
- All feed names should be correct.
- Keywords are not empty.
*/
package v1
