package testenvmanager

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/glog"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Event int

const (
	Group       = "testing.istio.io"
	Version     = "v1"
	EventAdd    = iota
	EventUpdate = iota
	EventDelete = iota
)

type EventHandler interface {
	Apply(obj interface{}, event Event) error
}

type CacheHandler struct {
	informer cache.SharedIndexInformer
	handler  EventHandler
}

type Controller struct {
	queue                           workqueue.RateLimitingInterface
	requestHandler, instanceHandler *CacheHandler
}

type Task struct {
	event   Event
	handler EventHandler
	obj     interface{}
}

func NewTask(handler EventHandler, obj interface{}, event Event) Task {
	return Task{handler: handler, obj: obj, event: event}
}

func NewController(
	rc, ic *crdclient,
	resyncPeriod time.Duration,
	requestHandler, instanceHandler EventHandler) *Controller {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "q1")

	c := &Controller{
		queue: queue,
	}

	c.requestHandler = &CacheHandler{
		informer: c.createInstanceInformer(rc.NewListWatch(), TestEnvRequestPlural, resyncPeriod, requestHandler),
		handler:  requestHandler,
	}

	c.instanceHandler = &CacheHandler{
		informer: c.createInstanceInformer(ic.NewListWatch(), TestEnvInstancePlural, resyncPeriod, instanceHandler),
		handler:  instanceHandler,
	}

	return c
}

func (c *Controller) createInstanceInformer(lw cache.ListerWatcher, kind string, resyncPeriod time.Duration, handler EventHandler) cache.SharedIndexInformer {
	informer := cache.NewSharedIndexInformer(
		lw,
		knownTypes[kind].object.DeepCopyObject(),
		resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{

			AddFunc: func(obj interface{}) {
				c.queue.Add(NewTask(handler, obj, EventAdd))
			},
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					c.queue.Add(NewTask(handler, cur, EventUpdate))
				}
			},
			DeleteFunc: func(obj interface{}) {
				c.queue.Add(NewTask(handler, obj, EventDelete))
			},
		},
	)
	return informer
}

// handleErr checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.queue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.queue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing  %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	utilruntime.HandleError(err)
	glog.Infof("Dropping %q out of the queue: %v", key, err)
}

// Run will start the controller.
// StopCh channel is used to send interrupt signal to stop it.
func (c *Controller) Run(stopCh <-chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer c.queue.ShutDown()

	glog.Info("Starting kubewatch controller")

	go c.requestHandler.informer.Run(stopCh)
	go c.instanceHandler.informer.Run(stopCh)

	// wait for the caches to synchronize before starting the worker
	if !cache.WaitForCacheSync(stopCh, c.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	glog.Info("Kubewatch controller synced and ready")

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) HasSynced() bool {
	if !c.instanceHandler.informer.HasSynced() {
		glog.V(2).Infof("Request Controller is syncing")
		return false
	}
	if !c.requestHandler.informer.HasSynced() {
		glog.V(2).Infof("Instance Controller is syncing")
		return false
	}
	return true
}

func (c *Controller) runWorker() {
	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	var err error
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	t, ok := key.(Task)
	if !ok {
		err = fmt.Errorf("Could not extract Handler from task %v", key)
	}

	err = t.handler.Apply(t.obj, t.event)
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}
