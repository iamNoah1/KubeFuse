package memory

import "k8s.io/client-go/discovery"

type MemCacheClient struct {
	discovery.RESTMapper
}

func NewMemCacheClient(d discovery.RESTMapper) *MemCacheClient {
	return &MemCacheClient{RESTMapper: d}
}
