package poolmanager

type ClusterProvider interface {
	Create(ClusterConfig) (ClusterInstance, error)
	Delete(ClusterInstance) error
	CheckStatus(id string) ClusterState
}
