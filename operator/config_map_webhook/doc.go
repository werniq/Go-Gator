// This package provides the implementation of a Kubernetes admission controller webhook
// designed to validate ConfigMap objects and trigger custom reconciliation logic for HotNews resources.
//
// Core Features:
// - Admission validation for ConfigMap objects.
// - Custom resource reconciliation for HotNews resources based on ConfigMap data.
//
// Key Functions:
//   - validatingConfigMapHandler: Handles HTTP requests from the Kubernetes admission controller and
//     delegates the validation logic.
//   - validateConfigMap: Ensures ConfigMap objects are well-formed and meet the required conditions.
//   - getAllHotNewsFromNamespace: Retrieves all HotNews resources from the specified namespace.
//   - triggerHotNewsReconcile: Triggers reconciliation of HotNews resources based on ConfigMap feed groups.
//   - isKubeNamespace: Determines if a namespace is a system namespace in Kubernetes.
package main
