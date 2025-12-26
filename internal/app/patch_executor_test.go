package app

import (
	"context"
	"testing"
	"time"

	"kubefuse/internal/domain"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func TestBuildMergePatchPayloadAddsAnnotations(t *testing.T) {
	intent := domain.NewPatchIntent(domain.ResourceRef{Kind: "deploy", Name: "web", Namespace: "prod"}, []domain.Patch{
		domain.NewPatch([]string{"spec", "replicas"}, domain.NewInt(3)),
		domain.NewPatch([]string{"metadata", "labels", "env"}, domain.NewString("prod")),
	},
		domain.WithReason("hotfix"),
		domain.WithTTL(5*time.Minute),
	)

	payload, err := buildMergePatchPayload(intent)
	if err != nil {
		t.Fatalf("buildMergePatchPayload returned error: %v", err)
	}

	metadataAny, ok := payload["metadata"].(map[string]any)
	if !ok {
		t.Fatalf("metadata map missing in payload")
	}

	annotations, ok := metadataAny["annotations"].(map[string]any)
	if !ok {
		t.Fatalf("annotations map missing in payload")
	}

	if annotations[annotationReason] != "hotfix" {
		t.Fatalf("expected reason annotation, got %v", annotations[annotationReason])
	}

	if annotations[annotationTTL] != "5m0s" {
		t.Fatalf("expected ttl annotation '5m0s', got %v", annotations[annotationTTL])
	}

	spec, ok := payload["spec"].(map[string]any)
	if !ok {
		t.Fatalf("expected spec map in payload")
	}

	if spec["replicas"] != int64(3) {
		t.Fatalf("expected replicas to be 3, got %v", spec["replicas"])
	}
}

func TestExecutePatchIntentAppliesPatch(t *testing.T) {
	gvk := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvrSingular := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployment"}

	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{{Group: "apps", Version: "v1"}})
	mapper.AddSpecific(gvk, gvr, gvrSingular, meta.RESTScopeNamespace)

	deployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name":      "web",
				"namespace": "default",
			},
			"spec": map[string]any{
				"replicas": int64(1),
			},
		},
	}

	scheme := runtime.NewScheme()
	client := fake.NewSimpleDynamicClient(scheme, deployment)
	executor := NewPatchExecutorWithClients(client, mapper)

	intent := domain.NewPatchIntent(domain.ResourceRef{Kind: "deployments", Name: "web", Namespace: "default"},
		[]domain.Patch{domain.NewPatch([]string{"spec", "replicas"}, domain.NewInt(4))},
		domain.WithReason("scale up"),
	)

	if err := executor.ExecutePatchIntent(context.Background(), intent); err != nil {
		t.Fatalf("ExecutePatchIntent returned error: %v", err)
	}

	updated, err := client.Resource(gvr).Namespace("default").Get(context.Background(), "web", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to fetch patched resource: %v", err)
	}

	replicas, found, err := unstructured.NestedInt64(updated.Object, "spec", "replicas")
	if err != nil {
		t.Fatalf("error reading replicas: %v", err)
	}

	if !found || replicas != 4 {
		t.Fatalf("expected replicas to be updated to 4, got %d (found=%t)", replicas, found)
	}

	reason, found, err := unstructured.NestedString(updated.Object, "metadata", "annotations", annotationReason)
	if err != nil {
		t.Fatalf("error reading reason annotation: %v", err)
	}

	if !found || reason != "scale up" {
		t.Fatalf("expected reason annotation to be set, got %q (found=%t)", reason, found)
	}
}

func TestExecutePatchIntentRollsBackAfterTTL(t *testing.T) {
	gvk := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	gvr := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	gvrSingular := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployment"}

	mapper := meta.NewDefaultRESTMapper([]schema.GroupVersion{{Group: "apps", Version: "v1"}})
	mapper.AddSpecific(gvk, gvr, gvrSingular, meta.RESTScopeNamespace)

	deployment := &unstructured.Unstructured{
		Object: map[string]any{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]any{
				"name":      "web",
				"namespace": "default",
				"annotations": map[string]any{
					annotationReason: "previous",
					annotationTTL:    "30s",
				},
			},
			"spec": map[string]any{
				"replicas": int64(1),
			},
		},
	}

	scheme := runtime.NewScheme()
	client := fake.NewSimpleDynamicClient(scheme, deployment)
	executor := NewPatchExecutorWithClients(client, mapper)
	executor.waitFunc = func(ctx context.Context, d time.Duration) error {
		return nil
	}

	intent := domain.NewPatchIntent(domain.ResourceRef{Kind: "deployments", Name: "web", Namespace: "default"},
		[]domain.Patch{domain.NewPatch([]string{"spec", "replicas"}, domain.NewInt(4))},
		domain.WithReason("scale up"),
		domain.WithTTL(5*time.Second),
	)

	if err := executor.ExecutePatchIntent(context.Background(), intent); err != nil {
		t.Fatalf("ExecutePatchIntent returned error: %v", err)
	}

	updated, err := client.Resource(gvr).Namespace("default").Get(context.Background(), "web", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to fetch patched resource: %v", err)
	}

	replicas, found, err := unstructured.NestedInt64(updated.Object, "spec", "replicas")
	if err != nil {
		t.Fatalf("error reading replicas: %v", err)
	}

	if !found || replicas != 1 {
		t.Fatalf("expected replicas to be rolled back to 1, got %d (found=%t)", replicas, found)
	}

	reason, found, err := unstructured.NestedString(updated.Object, "metadata", "annotations", annotationReason)
	if err != nil {
		t.Fatalf("error reading reason annotation: %v", err)
	}

	if !found || reason != "previous" {
		t.Fatalf("expected reason annotation to be restored, got %q (found=%t)", reason, found)
	}

	ttl, found, err := unstructured.NestedString(updated.Object, "metadata", "annotations", annotationTTL)
	if err != nil {
		t.Fatalf("error reading ttl annotation: %v", err)
	}

	if !found || ttl != "30s" {
		t.Fatalf("expected ttl annotation to be restored, got %q (found=%t)", ttl, found)
	}
}
