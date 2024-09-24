package v1

import "sigs.k8s.io/controller-runtime/pkg/client"

// SetupClient is used to initialize k8s client.
func SetupClient(client client.Client) {
	k8sClient = client
}
