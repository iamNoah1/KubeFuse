package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"kubefuse/internal/domain"
	"os"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	annotationReason = "kubefuse.dev/reason"
	annotationTTL    = "kubefuse.dev/ttl"
	fieldManager     = "kubefuse"
)

type PatchExecutor struct {
	dynamicClient dynamic.Interface
	restMapper    meta.RESTMapper
}

func NewPatchExecutor() (*PatchExecutor, error) {
	cfg, err := buildKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig: %w", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create discovery client: %w", err)
	}

	mapper := restmapper.NewShortcutExpander(
		restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(disco)),
		disco,
	)

	dynClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create dynamic client: %w", err)
	}

	return &PatchExecutor{dynamicClient: dynClient, restMapper: mapper}, nil
}

func NewPatchExecutorWithClients(dynamicClient dynamic.Interface, mapper meta.RESTMapper) *PatchExecutor {
	return &PatchExecutor{dynamicClient: dynamicClient, restMapper: mapper}
}

func (p *PatchExecutor) ExecutePatchIntent(ctx context.Context, intent domain.PatchIntent) error {
	if len(intent.Patches) == 0 {
		return errors.New("no patches to apply")
	}

	if p.dynamicClient == nil || p.restMapper == nil {
		return errors.New("executor is not configured with Kubernetes clients")
	}

	payload, err := buildMergePatchPayload(intent)
	if err != nil {
		return err
	}

	patchBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal patch payload: %w", err)
	}

	resource, err := p.resolveResource(intent.Resource.Kind)
	if err != nil {
		return err
	}

	resourceClient, err := p.resourceInterface(resource, intent.Resource.Namespace)
	if err != nil {
		return err
	}

	options := metav1.PatchOptions{FieldManager: fieldManager}
	if intent.DryRun {
		options.DryRun = []string{metav1.DryRunAll}
	}

	_, err = resourceClient.Patch(ctx, intent.Resource.Name, types.MergePatchType, patchBytes, options)
	if err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	return nil
}

func (p *PatchExecutor) resolveResource(kind string) (*meta.RESTMapping, error) {
	gvr, err := p.restMapper.ResourceFor(schema.GroupVersionResource{Resource: strings.ToLower(kind)})
	if err != nil {
		return nil, fmt.Errorf("unable to resolve resource for kind %q: %w", kind, err)
	}

	gvk, err := p.restMapper.KindFor(gvr)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve kind for resource %q: %w", gvr.Resource, err)
	}

	mapping, err := p.restMapper.RESTMapping(gvk.GroupKind(), gvr.Version)
	if err != nil {
		return nil, fmt.Errorf("unable to map resource %q: %w", gvr.Resource, err)
	}

	return mapping, nil
}

func (p *PatchExecutor) resourceInterface(mapping *meta.RESTMapping, namespace string) (dynamic.ResourceInterface, error) {
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		ns := namespace
		if ns == "" {
			ns = "default"
		}
		return p.dynamicClient.Resource(mapping.Resource).Namespace(ns), nil
	}

	return p.dynamicClient.Resource(mapping.Resource), nil
}

func buildMergePatchPayload(intent domain.PatchIntent) (map[string]any, error) {
	root := map[string]any{}

	for _, patch := range intent.Patches {
		if len(patch.Path) == 0 {
			return nil, errors.New("patch path cannot be empty")
		}

		current := root

		for i, segment := range patch.Path {
			if i == len(patch.Path)-1 {
				current[segment] = patch.Value.ToInterface()
				continue
			}

			next, ok := current[segment]
			if !ok {
				nested := map[string]any{}
				current[segment] = nested
				current = nested
				continue
			}

			nested, ok := next.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("conflicting patch path at %s", strings.Join(patch.Path[:i+1], "."))
			}

			current = nested
		}
	}

	if intent.Reason != "" || intent.TTL != 0 {
		metadata := ensureNestedMap(root, "metadata")
		annotations := ensureNestedMap(metadata, "annotations")

		if intent.Reason != "" {
			annotations[annotationReason] = intent.Reason
		}

		if intent.TTL != 0 {
			annotations[annotationTTL] = intent.TTL.String()
		}
	}

	return root, nil
}

func ensureNestedMap(parent map[string]any, key string) map[string]any {
	if val, ok := parent[key]; ok {
		if nested, ok := val.(map[string]any); ok {
			return nested
		}
	}

	nested := map[string]any{}
	parent[key] = nested
	return nested
}

func buildKubeConfig() (*rest.Config, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err == nil {
		return config, nil
	}

	inClusterConfig, inClusterErr := rest.InClusterConfig()
	if inClusterErr == nil {
		return inClusterConfig, nil
	}

	return nil, fmt.Errorf("kubeconfig error: %v, in-cluster error: %w", err, inClusterErr)
}
