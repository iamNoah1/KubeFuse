package dynamic

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceInterface interface {
	Namespace(string) ResourceInterface
	Patch(ctx context.Context, name string, pt interface{}, data []byte, opts metav1.PatchOptions) (*unstructured.Unstructured, error)
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*unstructured.Unstructured, error)
}

type Interface interface {
	Resource(schema.GroupVersionResource) ResourceInterface
}
