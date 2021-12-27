/*
Copyright 2021 The Clusternet Authors.

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

package description

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	appsapi "gaia.io/gaia/pkg/apis/apps/v1alpha1"
	gaiaClientSet "gaia.io/gaia/pkg/generated/clientset/versioned"
	appInformers "gaia.io/gaia/pkg/generated/informers/externalversions/apps/v1alpha1"
	platformInformers "gaia.io/gaia/pkg/generated/informers/externalversions/platform/v1alpha1"
	appListers "gaia.io/gaia/pkg/generated/listers/apps/v1alpha1"
)

// controllerKind contains the schema.GroupVersionKind for this controller type.
var controllerKind = appsapi.SchemeGroupVersion.WithKind("Description")

type SyncHandlerFunc func(description *appsapi.Description) error

// Controller is a controller that handle Description
type Controller struct {
	childGaiaClient gaiaClientSet.Interface

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface

	descLister      appListers.DescriptionLister
	descSynced      cache.InformerSynced
	mcSynced        cache.InformerSynced
	syncHandlerFunc SyncHandlerFunc
	cacheHandler    cache.ResourceEventHandler
}

func NewController(childGaiaClient gaiaClientSet.Interface, descInformer appInformers.DescriptionInformer,
	mcInformer platformInformers.ManagedClusterInformer, cacheHandler cache.ResourceEventHandler, syncHandlerFunc SyncHandlerFunc) (*Controller, error) {
	if syncHandlerFunc == nil || cacheHandler == nil {
		return nil, fmt.Errorf("syncHandlerFunc or cacheHandler must be set")
	}

	c := &Controller{
		childGaiaClient: childGaiaClient,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "description"),
		descLister:      descInformer.Lister(),
		descSynced:      descInformer.Informer().HasSynced,
		mcSynced: func() bool {
			return true
		},
		syncHandlerFunc: syncHandlerFunc,
	}

	// Manage the addition/update of Description
	descInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addDescription,
		UpdateFunc: c.updateDescription,
		DeleteFunc: c.deleteDescription,
	})

	if mcInformer != nil && cacheHandler != nil {
		mcInformer.Informer().AddEventHandler(cacheHandler)
		c.mcSynced = mcInformer.Informer().HasSynced
	}
	return c, nil
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("starting description controller...")
	defer klog.Info("shutting down description controller")

	// Wait for the caches to be synced before starting workers
	if !cache.WaitForNamedCacheSync("description-controller", stopCh, c.descSynced, c.mcSynced) {
		return
	}

	klog.V(5).Infof("starting %d worker threads", workers)
	// Launch workers to process Description resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
}

func (c *Controller) addDescription(obj interface{}) {
	desc := obj.(*appsapi.Description)
	klog.V(4).Infof("adding Description %q", klog.KObj(desc))
	c.Enqueue(desc)
}

func (c *Controller) updateDescription(old, cur interface{}) {
	oldDesc := old.(*appsapi.Description)
	newDesc := cur.(*appsapi.Description)

	if newDesc.DeletionTimestamp != nil {
		c.Enqueue(newDesc)
		return
	}

	// Decide whether discovery has reported a spec change.
	if reflect.DeepEqual(oldDesc.Spec, newDesc.Spec) {
		klog.V(4).Infof("no updates on the spec of Description %s, skipping syncing", klog.KObj(oldDesc))
		return
	}

	klog.V(4).Infof("updating Description %q", klog.KObj(oldDesc))
	c.Enqueue(newDesc)
}

func (c *Controller) deleteDescription(obj interface{}) {
	desc, ok := obj.(*appsapi.Description)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("couldn't get object from tombstone %#v", obj))
			return
		}
		desc, ok = tombstone.Obj.(*appsapi.Description)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("tombstone contained object that is not a Description %#v", obj))
			return
		}
	}
	klog.V(4).Infof("deleting Description %q", klog.KObj(desc))
	c.Enqueue(desc)
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Description resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("successfully synced Description %q", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Description resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// If an error occurs during handling, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.

	// Convert the namespace/name string into a distinct namespace and name
	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	klog.V(4).Infof("start processing Description %q", key)
	// Get the Description resource with this name
	desc, err := c.descLister.Descriptions(ns).Get(name)
	// The Description resource may no longer exist, in which case we stop processing.
	if errors.IsNotFound(err) {
		klog.V(2).Infof("Description %q has been deleted", key)
		return nil
	}
	if err != nil {
		return err
	}

	desc.Kind = controllerKind.Kind
	desc.APIVersion = controllerKind.Version
	err = c.syncHandlerFunc(desc)
	return err
}

func (c *Controller) UpdateDescriptionStatus(desc *appsapi.Description, status *appsapi.DescriptionStatus) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance

	klog.V(5).Infof("try to update Description %q status", desc.Name)

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		desc.Status = *status
		_, err := c.childGaiaClient.AppsV1alpha1().Descriptions(desc.Namespace).UpdateStatus(context.TODO(), desc, metav1.UpdateOptions{})
		if err == nil {
			//TODO
			return nil
		}

		if updated, err := c.descLister.Descriptions(desc.Namespace).Get(desc.Name); err == nil {
			// make a copy so we don't mutate the shared cache
			desc = updated.DeepCopy()
		} else {
			utilruntime.HandleError(fmt.Errorf("error getting updated Description %q from lister: %v", desc.Name, err))
		}
		return err
	})
}

// Enqueue takes a Description resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Description.
func (c *Controller) Enqueue(desc *appsapi.Description) {
	key, err := cache.MetaNamespaceKeyFunc(desc)
	if err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}
