package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	"istio.io/test-infra/toolbox/testenvmanager"
)

type fakeHandler struct {
}

func (h fakeHandler) Apply(o interface{}, e testenvmanager.Event) error {
	var n string
	r, ok := o.(*testenvmanager.TestEnvRequest)
	if !ok {
		i, ok := o.(*testenvmanager.TestEnvInstance)
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
		glog.Infof("Recycle env %v", n)
		return nil
	case testenvmanager.EventUpdate:
		glog.Infof("Update env %v", n)
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
	config, scheme, err := testenvmanager.CreateRESTConfig(kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}
	if err := testenvmanager.RegisterResources(config); err != nil {
		glog.Fatal(err)
	}

	// creates the clientset
	restClient, err := rest.RESTClientFor(config)
	if err != nil {
		glog.Fatal(err)
	}

	rc := testenvmanager.NewCRDClient(restClient, scheme, v1.NamespaceDefault, testenvmanager.TestEnvRequestPlural)
	ic := testenvmanager.NewCRDClient(restClient, scheme, v1.NamespaceDefault, testenvmanager.TestEnvInstancePlural)

	req := &testenvmanager.TestEnvRequest{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "testenv-req-",
		},
		Spec: testenvmanager.TestEnvRequestSpec{
			Config: testenvmanager.ClusterConfig{
				NumCoresPerNode: 4,
				TTL:             1 * time.Hour,
				RBAC:            true,
				NumNodes:        3,
			},
		},
	}

	rc.Create(req)

	rh := fakeHandler{}
	ih := fakeHandler{}

	controller := testenvmanager.NewController(rc, ic, 60*time.Second, rh, ih)

	// Now let's start the controller
	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(stop)

	// Wait forever
	select {}
}
