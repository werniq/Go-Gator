package controller

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// StatusUpdatePredicate
type StatusUpdatePredicate struct {
	predicate.Funcs
}

func (StatusUpdatePredicate) Update(e event.UpdateEvent) bool {
	oldObj := e.ObjectOld.DeepCopyObject().(metav1.Object)
	newObj := e.ObjectNew.DeepCopyObject().(metav1.Object)

	oldStatus := reflect.ValueOf(oldObj).Elem().FieldByName("Status").Interface()
	newStatus := reflect.ValueOf(newObj).Elem().FieldByName("Status").Interface()

	return !reflect.DeepEqual(oldStatus, newStatus)
}

func (StatusUpdatePredicate) Create(e event.CreateEvent) bool {
	return true
}

func (StatusUpdatePredicate) Delete(e event.DeleteEvent) bool {
	return true
}

func (StatusUpdatePredicate) Generic(e event.GenericEvent) bool {
	return true
}
