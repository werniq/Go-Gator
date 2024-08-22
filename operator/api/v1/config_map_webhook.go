package v1

//
//import (
//	v1 "k8s.io/api/core/v1"
//	"k8s.io/client-go/informers"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/cache"
//)
//
//// validateConfigMapWebhook validates the ConfigMap webhook by checking if the feeds with certain names exist
////
//// It retrieves the list of feed CRDs from the cluster and checks if the feeds from the ConfigMap exist
//func validateConfigMapWebhook() error {
//	var err error
//	clientset, err = kubernetes.NewForConfig(c)
//	if err != nil {
//		return err
//	}
//
//	factory := informers.NewSharedInformerFactory(clientset, 0)
//	informer := factory.Core().V1().ConfigMaps().Informer()
//
//	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
//		AddFunc: func(obj interface{}) {
//			configMap := obj.(*v1.ConfigMap)
//			err = checkIfFeedExists(configMap)
//			if err != nil {
//				hotnewslog.Error(err, "Failed to validate feeds")
//			}
//		},
//		UpdateFunc: func(oldObj, newObj interface{}) {
//			configMap := newObj.(*v1.ConfigMap)
//			err = checkIfFeedExists(configMap)
//			if err != nil {
//				hotnewslog.Error(err, "Failed to validate feeds")
//			}
//		},
//	})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func isInArray(feedList []Feed, keyword string) bool {
//	for _, feed := range feedList {
//		if feed.Name == keyword {
//			return true
//		}
//	}
//	return false
//}
