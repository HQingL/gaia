/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package scheduler

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"reflect"

	"k8s.io/client-go/tools/cache"
)

// Returns cache.ResourceEventHandlerFuncs that trigger the given function
// on all object changes.
func NewTriggerOnAllChanges(triggerFunc func(metav1.Object)) *cache.ResourceEventHandlerFuncs {
	return &cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(old interface{}) {
			if deleted, ok := old.(cache.DeletedFinalStateUnknown); ok {
				// This object might be stale but ok for our current usage.
				old = deleted.Obj
				if old == nil {
					return
				}
			}
			oldObj := old.(metav1.Object)
			triggerFunc(oldObj)
		},
		AddFunc: func(cur interface{}) {
			curObj := cur.(metav1.Object)
			triggerFunc(curObj)
		},
		UpdateFunc: func(old, cur interface{}) {
			curObj := cur.(metav1.Object)
			if !reflect.DeepEqual(old, cur) {
				triggerFunc(curObj)
			}
		},
	}
}