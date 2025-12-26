package app

import (
	"context"
	"fmt"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
)

func ListResourceKinds() ([]string, error) {
	cfg, err := buildKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig: %w", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create discovery client: %w", err)
	}

	resourceLists, err := disco.ServerPreferredResources()
	if err != nil && !discovery.IsGroupDiscoveryFailedError(err) {
		return nil, fmt.Errorf("unable to discover resources: %w", err)
	}

	seen := map[string]struct{}{}
	kinds := make([]string, 0)

	for _, list := range resourceLists {
		for _, res := range list.APIResources {
			if strings.Contains(res.Name, "/") {
				continue
			}

			if _, ok := seen[res.Name]; !ok {
				seen[res.Name] = struct{}{}
				kinds = append(kinds, res.Name)
			}

			for _, short := range res.ShortNames {
				if _, ok := seen[short]; ok {
					continue
				}
				seen[short] = struct{}{}
				kinds = append(kinds, short)
			}
		}
	}

	sort.Strings(kinds)
	return kinds, nil
}

func ListResourceNames(ctx context.Context, kind, namespace string) ([]string, error) {
	executor, err := NewPatchExecutor()
	if err != nil {
		return nil, err
	}

	mapping, err := executor.resolveResource(kind)
	if err != nil {
		return nil, err
	}

	resourceClient, err := executor.resourceInterface(mapping, namespace)
	if err != nil {
		return nil, err
	}

	list, err := resourceClient.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list %s resources: %w", kind, err)
	}

	names := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		names = append(names, item.GetName())
	}

	sort.Strings(names)
	return names, nil
}
