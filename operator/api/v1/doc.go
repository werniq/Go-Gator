/*
Package v1 contains CRD definitions for the Feed resource, and Validation webhook that ensures uniqueness of the
feed.

This package includes the following components:
- Custom Resource Definition for the Feed resource.
- Webhook for validating Feed resources.

The Feed resource represents a news feed with a name and a link.

The webhook ensures that the Feed resources meet the required validation criteria, in particular:
- Name should be unique within namespace, and should not be more than 20 symbols
- Endpoint of feed should be either http or https url.
*/
package v1
