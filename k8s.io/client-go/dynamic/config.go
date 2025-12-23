package dynamic

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type inMemoryClient struct {
	store map[string]*unstructured.Unstructured
	mu    sync.RWMutex
}

func key(gvr schema.GroupVersionResource, ns, name string) string {
	_ = gvr
	return ns + "/" + name
}

func NewForConfig(_ interface{}) (Interface, error) {
	return &inMemoryClient{store: map[string]*unstructured.Unstructured{}}, nil
}

func (c *inMemoryClient) Resource(gvr schema.GroupVersionResource) ResourceInterface {
	return &inMemoryResource{client: c, gvr: gvr}
}

type inMemoryResource struct {
	client *inMemoryClient
	gvr    schema.GroupVersionResource
	ns     string
}

func (r *inMemoryResource) Namespace(ns string) ResourceInterface {
	return &inMemoryResource{client: r.client, gvr: r.gvr, ns: ns}
}

func (r *inMemoryResource) Patch(ctx context.Context, name string, _ interface{}, data []byte, _ metav1.PatchOptions) (*unstructured.Unstructured, error) {
	r.client.mu.Lock()
	defer r.client.mu.Unlock()
	k := key(r.gvr, r.ns, name)
	current, ok := r.client.store[k]
	if !ok {
		current = &unstructured.Unstructured{Object: map[string]any{"metadata": map[string]any{"name": name, "namespace": r.ns}}}
		r.client.store[k] = current
	}
	var patch map[string]any
	if err := json.Unmarshal(data, &patch); err != nil {
		return nil, err
	}
	applyPatch(current.Object, patch)
	return current, nil
}

func (r *inMemoryResource) Get(ctx context.Context, name string, _ metav1.GetOptions) (*unstructured.Unstructured, error) {
	r.client.mu.RLock()
	defer r.client.mu.RUnlock()
	k := key(r.gvr, r.ns, name)
	if obj, ok := r.client.store[k]; ok {
		return obj, nil
	}
	return nil, fmt.Errorf("not found")
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
