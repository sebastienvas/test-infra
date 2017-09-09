package poolmanager

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ClusterRequestHandler interface {
	ProvisionCluster(r *ClusterRequest) error
	RecycleCluster(r *ClusterRequest) error
}

type Controller struct {
	dynamic  *rest.RESTClient
	queue    workqueue.RateLimitingInterface
	informer cache.SharedIndexInformer
	handler  ClusterRequestHandler
}

type OnDemand struct {
	cm ClusterManager
}

type FixedSizePools struct {
	cm        ClusterManager
	LifeSpan  time.Duration
	QueueSize int
}

type ClusterManagerMode interface {
	Get(ClusterConfig) (ClusterInstance, error)
	Recycle(ClusterInstance)
}

type ClusterManager struct {
	provider          ClusterProvider
	clustersInstances map[string][]ClusterInstance
}

func NewController(client *rest.RESTClient, namespace string, resyncPeriod time.Duration, handler ClusterRequestHandler) *Controller {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "q1")

	informer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (result runtime.Object, err error) {
				result = knownTypes[ClusterRequestsKind].collection.DeepCopyObject()
				err = client.Get().
					Namespace(namespace).
					Resource(ClusterRequestsKind).
					Do().
					Into(result)
				return
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				return client.Get().
					Prefix("watch").
					Namespace(namespace).
					Resource(ClusterRequestsKind).
					Watch()
			},
		},
		knownTypes[ClusterRequestsKind].object.DeepCopyObject(),
		resyncPeriod, cache.Indexers{})

	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{

			AddFunc: func(obj interface{}) {
				key, err := cache.MetaNamespaceIndexFunc(obj)
				if err == nil {
					queue.Add(key)
				}

			},
			DeleteFunc: func(obj interface{}) {
				key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
				if err == nil {
					queue.Add(key)
				}
			},
		},
	)

	return &Controller{
		dynamic:  client,
		queue:    queue,
		informer: informer,
		handler:  handler,
	}
}

func (c *Controller) processItem(key string) error {
	glog.Infof("Processing change to Pod %s", key)
	var err error
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", key, err)
	}

	if !exists {
		err = c.handler.RecycleCluster(obj.(*ClusterRequest))

	} else {
		err = c.handler.ProvisionCluster(obj.(*ClusterRequest))
	}

	c.handleErr(err, key)
	return nil
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
		glog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)
	utilruntime.HandleError(err)
	glog.Infof("Dropping pod %q out of the queue: %v", key, err)
}

// Run will start the controller.
// StopCh channel is used to send interrupt signal to stop it.
func (c *Controller) Run(stopCh <-chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer c.queue.ShutDown()

	glog.Info("Starting kubewatch controller")

	go c.informer.Run(stopCh)

	// wait for the caches to synchronize before starting the worker
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	glog.Info("Kubewatch controller synced and ready")

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
	// processNextWorkItem will automatically wait until there's work available
	for c.processNextItem() {
		// continue looping
	}
}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.queue.Done(key)

	// Invoke the method containing the business logic
	err := c.processItem(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}
