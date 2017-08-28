package poolmanager

import "time"

type ClusterStatus string

const (
	CREATING = ClusterStatus("CREATING")
	READY    = ClusterStatus("READY")
	IN_USE   = ClusterStatus("IN_USE")
	DELETED  = ClusterStatus("DELETED ")
)

type ClusterConfig struct {
	TTL             time.Duration
	NumCoresPerNode int
	NumNodes        int
	Provider        string
	Version         string
	RBAC            bool
}

type ClusterRequest struct {
	Config ClusterConfig
	KubeConfig string
}

type ClusterInstance struct {
	Config     ClusterConfig
	Status     ClusterStatus
	ID         string
	KubeConfig string
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

