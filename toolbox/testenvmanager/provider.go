package testenvmanager

type ClusterProvider interface {
	Create(ClusterConfig) (TestEnvInstance, error)
	Delete(TestEnvInstance) error
	CheckStatus(id string) TestEnvState
}
