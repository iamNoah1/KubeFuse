package rest

type Config struct{}

func InClusterConfig() (*Config, error) {
	return nil, ErrNotInCluster
}

type inClusterError string

func (e inClusterError) Error() string { return string(e) }

const ErrNotInCluster inClusterError = "not in cluster"
