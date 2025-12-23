package fake

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type SimpleDynamicClient struct {
	store map[string]*unstructured.Unstructured
	mu    sync.RWMutex
}

func NewSimpleDynamicClient(_ interface{}, objs ...*unstructured.Unstructured) *SimpleDynamicClient {
	c := &SimpleDynamicClient{store: map[string]*unstructured.Unstructured{}}
	for _, obj := range objs {
		metaAny, _ := obj.Object["metadata"].(map[string]any)
		ns, _ := metaAny["namespace"].(string)
		name, _ := metaAny["name"].(string)
		key := c.key(schema.GroupVersionResource{}, ns, name)
		c.store[key] = obj
	}
	return c
}

func (c *SimpleDynamicClient) key(gvr schema.GroupVersionResource, ns, name string) string {
	_ = gvr
	return ns + "/" + name
}

func (c *SimpleDynamicClient) Resource(gvr schema.GroupVersionResource) dynamic.ResourceInterface {
	return &resourceClient{client: c, gvr: gvr}
}

type resourceClient struct {
	client *SimpleDynamicClient
	gvr    schema.GroupVersionResource
	ns     string
}

func (r *resourceClient) Namespace(ns string) dynamic.ResourceInterface {
	return &resourceClient{client: r.client, gvr: r.gvr, ns: ns}
}

func (r *resourceClient) Patch(ctx context.Context, name string, _ interface{}, data []byte, _ metav1.PatchOptions) (*unstructured.Unstructured, error) {
	r.client.mu.Lock()
	defer r.client.mu.Unlock()

	key := r.client.key(r.gvr, r.ns, name)
	current, ok := r.client.store[key]
	if !ok {
		return nil, fmt.Errorf("resource not found")
	}

	var patch map[string]any
	if err := json.Unmarshal(data, &patch); err != nil {
		return nil, err
	}

	applyPatch(current.Object, patch)
	return current, nil
}

func applyPatch(target map[string]any, patch map[string]any) {
	for k, v := range patch {
		if subPatch, ok := v.(map[string]any); ok {
			if existing, ok := target[k].(map[string]any); ok {
				applyPatch(existing, subPatch)
				continue
			}
			target[k] = copyMap(subPatch)
			continue
		}
		target[k] = v
	}
}

func copyMap(src map[string]any) map[string]any {
	dst := map[string]any{}
	for k, v := range src {
		if nested, ok := v.(map[string]any); ok {
			dst[k] = copyMap(nested)
			continue
		}
		dst[k] = v
	}
	return dst
}

func (r *resourceClient) Get(ctx context.Context, name string, _ metav1.GetOptions) (*unstructured.Unstructured, error) {
	r.client.mu.RLock()
	defer r.client.mu.RUnlock()
	key := r.client.key(r.gvr, r.ns, name)
	if obj, ok := r.client.store[key]; ok {
		return obj, nil
	}
	return nil, fmt.Errorf("not found")
}
