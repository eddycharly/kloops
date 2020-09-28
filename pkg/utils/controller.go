package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	toolscache "k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/cache"
)

func NewController(ch chan<- SocketData, informer cache.Informer, onCreated, onUpdated, onDeleted string) {
	informer.AddEventHandler(toolscache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			data := SocketData{
				MessageType: onCreated,
				Payload:     obj,
			}
			ch <- data
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldMeta, newMeta := oldObj.(metav1.Object), newObj.(metav1.Object)
			// If resourceVersion differs between old and new, an actual update event was observed
			if oldMeta.GetResourceVersion() != newMeta.GetResourceVersion() {
				data := SocketData{
					MessageType: onUpdated,
					Payload:     newObj,
				}
				ch <- data
			}
		},
		DeleteFunc: func(obj interface{}) {
			data := SocketData{
				MessageType: onDeleted,
				Payload:     obj,
			}
			ch <- data
		},
	})
}
