### Operator Overview
The `FeedReconciler` handles the creation, updating, and deletion of feeds on a news-aggregator server.
When a news feed is registered, a request is sent to the server to create a feed
with the specified name and endpoint.
Similarly, when a news feed is updated or deleted, requests are made to update
or remove the feed from the server accordingly.

The `HotNewsReconciler` is a Kubernetes controller responsible for managing the `HotNews` custom resource. 
It ensures that the status of the `HotNews` resource is always up-to-date by interacting with a news aggregator server
to fetch the latest news based on specified parameters.

### Description

The `FeedReconciler` performs the management of feeds on a news-aggregator server by creating,
updating, or deleting feeds.
When the status of a Feed CRD changes (whether it is created, updated, or deleted),
the operator sends a request to the news-aggregator server to perform the corresponding action
with the specified name and endpoint.

The `HotNewsReconciler` performs the following key tasks:
- Monitors the status of `HotNews` custom resources and triggers updates whenever the resource is created, updated, or deleted.
- Sends requests to a news aggregator server to retrieve news articles based on specified keywords, date ranges, and sources.
- Validates the input parameters (keywords, date range, and feeds) before making a request to ensure accuracy.
- Watches for changes in the `ConfigMap` containing feed groups and in the `Feed` CRD, and updates the `HotNews` resource accordingly.
- Handles the creation, update, and deletion of `HotNews` resources, ensuring that the Kubernetes cluster state aligns with the desired state.

The reconciler requeues and reconciles every 24 hours to keep the news content fresh and updated.

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```
