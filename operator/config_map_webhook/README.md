# ConfigMap Admission Controller Webhook
## Overview

This project implements an **Admission Controller Webhook** for Kubernetes that validates 
`ConfigMap` objects before they are admitted into the cluster. 
It ensures that every time a feed in the `ConfigMap` is updated, all `HotNews` resources that
use the same feed group are reconciled.

## Features

- Validates `ConfigMap` objects in Kubernetes.
- Ensures `ConfigMap` has a data field.
- Triggers reconciliation of `HotNews` resources based on the `ConfigMap` content.
- Handles requests through an HTTP webhook.

## Key Components

### `validatingConfigMapHandler`

- Handles incoming HTTP requests from Kubernetes for admission control.
- Parses the request and delegates validation logic to the `validateConfigMap` function.
- Returns an appropriate admission response, allowing or denying the request.

### `validateConfigMap`

- Validates that the `ConfigMap` has the required `data` field.
- Fetches all `HotNews` resources from the same namespace.
- Triggers a reconcile process for `HotNews` that contain the same `Feed` groups as in the `ConfigMap`.

### `getAllHotNewsFromNamespace`

- Retrieves all `HotNews` custom resources from the specified namespace.

### `triggerHotNewsReconcile`

- Triggers a reconcile for all `HotNews` resources that match feed groups from the `ConfigMap`.

### `isKubeNamespace`

- Helper function that checks if a given namespace is a Kubernetes-owned namespace.
