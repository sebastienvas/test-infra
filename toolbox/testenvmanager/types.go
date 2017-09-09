package testenvmanager

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	TestEnvRequestsKind   = "TestEnvRequest"
	TestEnvRequestPlural  = "testenvrequests"
	TestEnvInstancesKind  = "TestEnvInstance"
	TestEnvInstancePlural = "testenvinstances"
)

type TestEnvState string

var knownTypes = map[string]struct {
	object     runtime.Object
	collection runtime.Object
}{
	TestEnvRequestPlural: {
		object:     &TestEnvRequest{},
		collection: &TestEnvRequestList{},
	},
	TestEnvInstancePlural: {
		object:     &TestEnvInstance{},
		collection: &TestEnvInstanceList{},
	},
}

const (
	CREATING = TestEnvState("CREATING")
	READY    = TestEnvState("READY")
	IN_USE   = TestEnvState("IN_USE")
	DELETED  = TestEnvState("DELETED ")
)

type TestEnvRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TestEnvRequestSpec `json:"spec,omitempty"`
	Status            *TestEnvState      `json:"status,omitempty"`
}

type TestEnvRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*TestEnvRequest `json:"items"`
}

type ClusterConfig struct {
	TTL             time.Duration `json:"ttl,omitempty"`
	NumCoresPerNode int           `json:"numCoresPerNode,omitempty"`
	NumNodes        int           `json:"numNodes,omitempty"`
	Version         string        `json:"version,omitempty"`
	RBAC            bool          `json:"rbac,omitempty"`
}

type TestEnvRequestSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
}

type TestEnvStatus struct {
	State      TestEnvState `json:"state,omitempty"`
	KubeConfig string       `json:"kubeConfig,omitempty"`
}

type TestEnvInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TestEnvInstanceSpec `json:"spec,omitempty"`
	Status            *TestEnvState       `json:"status,omitempty"`
}

type TestEnvInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*TestEnvInstance `json:"items"`
}

type TestEnvInstanceSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
	ID     string
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
