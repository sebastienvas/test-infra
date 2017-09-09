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
	ClusterRequestName   = "clusterinstances"
)

type ClusterStatus string

const (
	CREATING = ClusterStatus("CREATING")
	READY    = ClusterStatus("READY")
	IN_USE   = ClusterStatus("IN_USE")
	DELETED  = ClusterStatus("DELETED ")
)

type ClusterRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterRequestSpec `json:"spec,omitempty"`
	Status            *ClusterStatus     `json:"status,omitempty"`
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
	State      ClusterStatus `json:"state,omitempty"`
	KubeConfig string        `json:"kubeConfig,omitempty"`
}

type ClusterInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterInstanceSpec `json:"spec,omitempty"`
	Status            *ClusterStatus      `json:"status,omitempty"`
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

func (r *ClusterRequest) DeepCopyObject() runtime.Object {
	//TODO: implement
	return nil
}
