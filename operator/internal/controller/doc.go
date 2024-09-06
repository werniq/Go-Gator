/*
Package controller contains the controller logic for the operator.

	It is a K8S extension that will help us to manage the Feed and HotNews resources.

	We have defined a CRDs called Feed and HotNews in the api package. The controller will watch for the Feed and HotNews
	resources and perform actions based on the status of the Feed resources.

	For example, when a new Feed is created, the controller will retrieve that resource, and make a request to
	out Go-Gator server to create a new source there.

	When a Feed is deleted, the controller will make a request to the Go-Gator server to delete the source.
	Same with the update operation.

	FeedReceiver is the main controller that will watch for the Feed resources.
	It has few methods:
	- Reconcile: this method will be called when a new resource is created, updated or deleted.
	- SetupWithManager: configure the controller with the manager
	- handleCreate: create a new source on the Go-Gator server
	- processHotNews: update a source on the Go-Gator server
	- handleDelete: delete a source on the Go-Gator server
	- initFeedStatus: initializes custom status fields for the Feed resource
	- updateFeedStatus: updates the status of the Feed resource

	HotNews is a CRD that will be used to retrieve news by the criteria, specified in the HotNewsSpec.
	For example, we can specify keywords, date range, feeds and feed groups.
	And then we will make requests to our news aggregator server with these parameters, and get the news.

	When the HotNews object is created, the controller will make a request to the Go-Gator server to get the news.

	HotNewsReceiver is the controller that will watch for the HotNews resources.
	It has few methods:
	- Reconcile: this method will be called when a new resource is created, updated or deleted.
	- SetupWithManager: configure the controller with the manager
	- handleCreate: create a new hot news group on the Go-Gator server
	- processHotNews: update a hot news group on the Go-Gator server
	- handleDelete: delete a hot news group on the Go-Gator server
	- initHotNewsStatus: initializes custom status fields for the HotNews resource
	- constructRequestUrl: constructs a request URL to get the news with the specified arguments
	- processFeedGroups and getFeedGroups are used to get data from feed-group-source config map, and process it
*/
package controller
