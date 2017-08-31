package poolmanager

import "time"

type ClusterStatus string

const (
	CREATING = ClusterStatus("CREATING")
	READY    = ClusterStatus("READY")
	IN_USE   = ClusterStatus("IN_USE")
	DELETED  = ClusterStatus("DELETED ")
)

type ObjectMeta struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	ResourceVersion string `json:"resourceVersion,omitempty"`
	UID             string `json:"uid,omitempty"`
}

type ClusterRequest struct {
	APIVersion string             `json:"apiVersion,omitempty"`
	Kind       string             `json:"kind,omitempty"`
	Metadata   ObjectMeta         `json:"metadata,omitempty"`
	Spec       ClusterRequestSpec `json:"spec,omitempty"`
	Status     ClusterStatus      `json:"status,omitempty"`
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
	APIVersion string              `json:"apiVersion,omitempty"`
	Kind       string              `json:"kind,omitempty"`
	Metadata   ObjectMeta          `json:"metadata,omitempty"`
	Spec       ClusterInstanceSpec `json:"spec,omitempty"`
	Status     ClusterStatus       `json:"status,omitempty"`
}

type ClusterInstanceSpec struct {
	Config ClusterConfig `json:"config,omitempty"`
	ID     string
}

type ClusterProvider interface {
	Create(ClusterConfig) (ClusterInstance, error)
	Delete(ClusterInstance) error
	CheckStatus(id string) ClusterStatus
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
	provider ClusterProvider
	Clusters map[string][]ClusterInstance
}

func (c *ClusterManager) Persist(*ClusterInstance) error
func (c *ClusterManager) Restore() error
func (c *ClusterManager) Create(*ClusterConfig) error
func (c *ClusterManager) Delete(*ClusterInstance) error
func (c *ClusterManager) Get(*ClusterConfig) error
func (c *ClusterManager) List(*ClusterConfig) ([]ClusterInstance, error)
