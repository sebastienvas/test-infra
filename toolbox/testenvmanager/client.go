package testenvmanager

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateRESTConfig for cluster API server, pass empty config file for in-cluster
func CreateRESTConfig(kubeconfig string) (config *rest.Config, types *runtime.Scheme, err error) {
	if kubeconfig == "" {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return
	}

	version := schema.GroupVersion{
		Group:   Group,
		Version: Version,
	}

	config.GroupVersion = &version
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON

	types = runtime.NewScheme()
	schemeBuilder := runtime.NewSchemeBuilder(
		func(scheme *runtime.Scheme) error {
			for _, kind := range knownTypes {
				scheme.AddKnownTypes(version, kind.object, kind.collection)
			}
			meta_v1.AddToGroupVersion(scheme, version)
			return nil
		})
	err = schemeBuilder.AddToScheme(types)
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(types)}

	return
}

// RegisterResources sends a request to create CRDs and waits for them to initialize
func RegisterResources(config *rest.Config) error {
	c, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return err
	}

	for _, s := range []struct{ p, k string }{
		{
			p: TestEnvRequestPlural,
			k: TestEnvRequestsKind,
		},
		{
			p: TestEnvInstancePlural,
			k: TestEnvInstancesKind,
		}} {
		crd := &apiextensionsv1beta1.CustomResourceDefinition{
			ObjectMeta: meta_v1.ObjectMeta{
				Name: fmt.Sprintf("%s.%s", s.p, Group),
			},
			Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
				Group:   Group,
				Version: Version,
				Scope:   apiextensionsv1beta1.NamespaceScoped,
				Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
					Plural: s.p,
					Kind:   s.k,
				},
			},
		}
		if _, err := c.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func NewCRDClient(cl *rest.RESTClient, scheme *runtime.Scheme, namespace, plural string) *crdclient {
	return &crdclient{cl: cl, ns: namespace, plural: plural,
		codec: runtime.NewParameterCodec(scheme)}
}

type crdclient struct {
	cl     *rest.RESTClient
	ns     string
	plural string
	codec  runtime.ParameterCodec
}

func (c *crdclient) Create(obj runtime.Object) (runtime.Object, error) {
	result := knownTypes[c.plural].object.DeepCopyObject()
	err := c.cl.Post().
		Namespace(c.ns).Resource(c.plural).
		Body(obj).Do().Into(result)
	return result, err
}

func (c *crdclient) Update(obj runtime.Object) (runtime.Object, error) {
	result := knownTypes[c.plural].object.DeepCopyObject()
	err := c.cl.Put().
		Namespace(c.ns).Resource(c.plural).
		Body(obj).Do().Into(result)
	return result, err
}

func (c *crdclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.cl.Delete().
		Namespace(c.ns).Resource(c.plural).
		Name(name).Body(options).Do().
		Error()
}

func (c *crdclient) Get(name string) (runtime.Object, error) {
	result := knownTypes[c.plural].object.DeepCopyObject()
	err := c.cl.Get().
		Namespace(c.ns).Resource(c.plural).
		Name(name).Do().Into(result)
	return result, err
}

func (c *crdclient) List(opts meta_v1.ListOptions) (runtime.Object, error) {
	result := knownTypes[c.plural].collection.DeepCopyObject()
	err := c.cl.Get().
		Namespace(c.ns).Resource(c.plural).
		VersionedParams(&opts, c.codec).
		Do().Into(result)
	return result, err
}

// Create a new List watch for our TPR
func (c *crdclient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(c.cl, c.plural, c.ns, fields.Everything())
}
