package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	"istio.io/test-infra/toolbox/testenvmanager"
)

type fakeHandler struct {
}

func (h fakeHandler) Apply(o interface{}, e testenvmanager.Event) error {
	var n string
	r, ok := o.(testenvmanager.TestEnvRequest)
	if !ok {
		i, ok := o.(testenvmanager.TestEnvInstance)
		if !ok {
			return fmt.Errorf("cannot construct request from %v", o)
		} else {
			n = i.Name
		}

	} else {
		n = r.Name
	}
	switch e {
	case testenvmanager.EventAdd:
		glog.Infof("Created Request %s", n)
		return nil
	case testenvmanager.EventDelete:
		glog.Infof("Recycle Cluser %v", n)
		return nil
	default:
		return fmt.Errorf("unkown event")
	}
}

func main() {
	var kubeconfig string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.Parse()

	// creates the connection
	config, err := testenvmanager.CreateRESTConfig(kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}

	// creates the clientset
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		glog.Fatal(err)
	}

	rh := fakeHandler{}
	ih := fakeHandler{}

	controller := testenvmanager.NewController(restClient, v1.NamespaceDefault, 60*time.Second, rh, ih)

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}
