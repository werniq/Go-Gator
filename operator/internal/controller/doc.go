// Package controller contains the controller logic for the operator.
// It is a K8S extension that will help us to manage the Feed resources.
//
// We have defined a CRD called Feed in the api package. The controller will watch for the Feed resources and
// perform actions based on the status of the Feed resources.
//
// For example, when a new resource is created, the controller will retrieve that resource, and make a request to
// out Go-Gator server to create a new source there.
//
// When a resource is deleted, the controller will make a request to the Go-Gator server to delete the source.
// Same with the update operation.
//
// FeedReceiver is the main controller that will watch for the Feed resources.
// It has few methods:
// - Reconcile: this method will be called when a new resource is created, updated or deleted.
// - SetupWithManager: configure the controller with the manager
// - handleCreate: create a new source on the Go-Gator server
// - handleUpdate: update a source on the Go-Gator server
// - handleDelete: delete a source on the Go-Gator server
// - initFeedStatus: initializes custom status fields for the Feed resource
// - updateFeedStatus: updates the status of the Feed resource
package controller
