package testenvmanager

import (
	"fmt"

	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateRESTConfig for cluster API server, pass empty config file for in-cluster
func CreateRESTConfig(kubeconfig string) (config *rest.Config, err error) {
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

	types := runtime.NewScheme()
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

func (in *TestEnvRequestList) DeepCopyInto(out *TestEnvRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = in.Items
}

func (in *TestEnvRequestList) DeepCopy() *TestEnvRequestList {
	if in == nil {
		return nil
	}
	out := new(TestEnvRequestList)
	in.DeepCopyInto(out)
	return out
}

func (in *TestEnvRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TestEnvRequest) DeepCopyInto(out *TestEnvRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *TestEnvRequest) DeepCopy() *TestEnvRequest {
	if in == nil {
		return nil
	}
	out := new(TestEnvRequest)
	in.DeepCopyInto(out)
	return out
}

func (in *TestEnvRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *TestEnvInstanceList) DeepCopyInto(out *TestEnvInstanceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = in.Items
}

func (in *TestEnvInstanceList) DeepCopy() *TestEnvInstanceList {
	if in == nil {
		return nil
	}
	out := new(TestEnvInstanceList)
	in.DeepCopyInto(out)
	return out
}

func (in *TestEnvInstanceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
func (in *TestEnvInstance) DeepCopyInto(out *TestEnvInstance) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *TestEnvInstance) DeepCopy() *TestEnvInstance {
	if in == nil {
		return nil
	}
	out := new(TestEnvInstance)
	in.DeepCopyInto(out)
	return out
}

func (in *TestEnvInstance) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
