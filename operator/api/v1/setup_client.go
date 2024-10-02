package v1

import "sigs.k8s.io/controller-runtime/pkg/client"

// k8sClient is a kubernetes client that is used to interact with the k8s API
var k8sClient client.Client

// SetupClient is used to initialize k8s client.
func SetupClient(client client.Client) {
	k8sClient = client
}
