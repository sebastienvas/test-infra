package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	"istio.io/test-infra/toolbox/poolmanager"
)

type fakeHandler struct {
}

func (h fakeHandler) ProvisionCluster(r *poolmanager.ClusterRequest) error {
	glog.Infof("Created Request %v", r.Name)
	return nil
}

func (h fakeHandler) RecycleCluster(r *poolmanager.ClusterRequest) error {
	glog.Infof("RecycleCluser %v", r.Name)
	return nil
}

func main() {
	var kubeconfig string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	// creates the connection
	config, err := poolmanager.CreateRESTConfig(kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}

	// creates the clientset
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		glog.Fatal(err)
	}

	handler := fakeHandler{}

	controller := poolmanager.NewController(restClient, v1.NamespaceDefault, 60*time.Second, handler)

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}
