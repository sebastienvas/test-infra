package provider

import "istio.io/test-infra/toolbox/clusterpool"

type ClusterProvider interface {
	Create(poolmanager.ClusterConfig) (poolmanager.ClusterInstance, error)
	Delete(poolmanager.ClusterInstance) error
	CheckStatus(id string) poolmanager.ClusterStatus
}
