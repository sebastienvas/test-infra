package poolmanager

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	ClusterRequestsKind  = "ClusterRequest"
	ClusterRequestName   = "clusterrequests"
	ClusterInstancesKind = "ClusterInstance"
	ClusterInstanceName  = "clusterinstances"
)

type ClusterState string

var knownTypes = map[string]struct {
	object     runtime.Object
	collection runtime.Object
}{
	ClusterRequestsKind: {
		object:     &ClusterRequest{},
		collection: &ClusterRequestList{},
	},
}

const (
	CREATING = ClusterState("CREATING")
	READY    = ClusterState("READY")
	IN_USE   = ClusterState("IN_USE")
	DELETED  = ClusterState("DELETED ")
)

type ClusterRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterRequestSpec `json:"spec,omitempty"`
	Status            *ClusterState      `json:"status,omitempty"`
}

type ClusterRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*ClusterRequest `json:"items"`
}

type ClusterConfig struct {
	TTL             time.Duration `json:"ttl,omitempty"`
	NumCoresPerNode int           `json:"numCoresPerNode,omitempty"`
	NumNodes        int           `json:"numNodes,omitempty"`
	Version         string        `json:"version,omitempty"`
	RBAC            bool          `json:"rbac,omitempty"`
}

type ClusterRequestSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
}

type ClusterStatus struct {
	State      ClusterState `json:"state,omitempty"`
	KubeConfig string       `json:"kubeConfig,omitempty"`
}

type ClusterInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterInstanceSpec `json:"spec,omitempty"`
	Status            *ClusterState       `json:"status,omitempty"`
}

type ClusterInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []*ClusterRequest `json:"items"`
}

type ClusterInstanceSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
	ID     string
}

func (in *ClusterRequestList) DeepCopyInto(out *ClusterRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	out.Items = in.Items
}

func (in *ClusterRequestList) DeepCopy() *ClusterRequestList {
	if in == nil {
		return nil
	}
	out := new(ClusterRequestList)
	in.DeepCopyInto(out)
	return out
}

func (in *ClusterRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *ClusterRequest) DeepCopyInto(out *ClusterRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *ClusterRequest) DeepCopy() *ClusterRequest {
	if in == nil {
		return nil
	}
	out := new(ClusterRequest)
	in.DeepCopyInto(out)
	return out
}

func (in *ClusterRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}
