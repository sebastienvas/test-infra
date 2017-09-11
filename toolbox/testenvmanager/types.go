package testenvmanager

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	TestEnvRequestsKind  = "ClusterRequest"
	TestEnvRequestName   = "clusterrequests"
	TestEnvInstancesKind = "ClusterInstance"
	TestEnvInstanceName  = "clusterinstances"
)

type TestEnvState string

var knownTypes = map[string]struct {
	object     runtime.Object
	collection runtime.Object
}{
	TestEnvRequestsKind: {
		object:     &TestEnvRequest{},
		collection: &TestEnvRequestList{},
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
	Items           []*TestEnvRequest `json:"items"`
}

type TestEnvInstanceSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
	ID     string
}
