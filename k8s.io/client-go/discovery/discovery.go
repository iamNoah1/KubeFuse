package discovery

import "k8s.io/apimachinery/pkg/runtime/schema"

type RESTMapper interface {
	ResourceFor(schema.GroupVersionResource) (schema.GroupVersionResource, error)
}

type DiscoveryClient struct{}

func NewDiscoveryClientForConfig(_ interface{}) (*DiscoveryClient, error) {
	return &DiscoveryClient{}, nil
}

func (DiscoveryClient) ResourceFor(res schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	return res, nil
}
