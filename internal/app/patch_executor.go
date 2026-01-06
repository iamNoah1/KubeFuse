package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/iamNoah1/KubeFuse/internal/domain"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
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
	spinnerInterval  = 120 * time.Millisecond
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

type PatchExecutor struct {
	dynamicClient dynamic.Interface
	restMapper    meta.RESTMapper
	waitFunc      func(context.Context, time.Duration) error
	output        io.Writer
}

func (p *PatchExecutor) ensureClients() error {
	if p.dynamicClient != nil && p.restMapper != nil {
		return nil
	}

	cfg, err := buildKubeConfig()
	if err != nil {
		return fmt.Errorf("unable to load kubeconfig: %w", err)
	}

	disco, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return fmt.Errorf("unable to create discovery client: %w", err)
	}

	mapper := restmapper.NewShortcutExpander(
		restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(disco)),
		disco,
		nil,
	)

	dynClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return fmt.Errorf("unable to create dynamic client: %w", err)
	}

	p.dynamicClient = dynClient
	p.restMapper = mapper
	return nil
}

func NewPatchExecutor() (*PatchExecutor, error) {
	return &PatchExecutor{waitFunc: waitForTTL, output: os.Stdout}, nil
}

func NewPatchExecutorWithClients(dynamicClient dynamic.Interface, mapper meta.RESTMapper) *PatchExecutor {
	return &PatchExecutor{dynamicClient: dynamicClient, restMapper: mapper, waitFunc: waitForTTL, output: os.Stdout}
}

func (p *PatchExecutor) ExecutePatchIntent(ctx context.Context, intent domain.PatchIntent) (err error) {
	if len(intent.Patches) == 0 {
		return errors.New("no patches to apply")
	}

	var stopSpinner func()
	if intent.TTL > 0 && !intent.DryRun {
		stopSpinner = p.startSpinner("Applying patch and scheduling rollback...")
		defer func() {
			if stopSpinner == nil {
				return
			}
			stopSpinner()
			if err != nil {
				writer := p.output
				if writer == nil {
					writer = os.Stdout
				}
				fmt.Fprintln(writer)
			}
		}()
	}

	if err := p.ensureClients(); err != nil {
		return err
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

	var current *unstructured.Unstructured
	if intent.DryRun {
		var err error
		current, err = resourceClient.Get(ctx, intent.Resource.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to read resource for preview: %w", err)
		}
	} else if intent.TTL > 0 {
		var err error
		current, err = resourceClient.Get(ctx, intent.Resource.Name, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("failed to read resource for rollback: %w", err)
		}
	}

	var rollbackPayload map[string]any
	if intent.TTL > 0 {
		var err error
		rollbackPayload, err = buildRollbackPayload(current, intent)
		if err != nil {
			return err
		}
	}

	if intent.DryRun {
		return p.printDryRunPreview(intent, payload, rollbackPayload)
	}

	options := metav1.PatchOptions{FieldManager: fieldManager}

	_, err = resourceClient.Patch(ctx, intent.Resource.Name, types.MergePatchType, patchBytes, options)
	if err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	if intent.DryRun || intent.TTL <= 0 {
		return nil
	}

	if p.waitFunc == nil {
		return errors.New("no wait function configured for TTL rollback")
	}

	if stopSpinner != nil {
		stopSpinner()
		stopSpinner = nil
	}

	if err := p.waitForRollback(ctx, intent.TTL); err != nil {
		return err
	}

	rollbackBytes, err := json.Marshal(rollbackPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal rollback payload: %w", err)
	}

	_, err = resourceClient.Patch(ctx, intent.Resource.Name, types.MergePatchType, rollbackBytes, options)
	if err != nil {
		return fmt.Errorf("failed to rollback patch: %w", err)
	}

	return nil
}

func (p *PatchExecutor) startSpinner(message string) func() {
	writer := p.output
	if writer == nil {
		writer = os.Stdout
	}

	done := make(chan struct{})
	stopped := make(chan struct{})
	ticker := time.NewTicker(spinnerInterval)
	step := 0

	render := func() {
		fmt.Fprintf(writer, "\r%s %s", spinnerFrames[step%len(spinnerFrames)], message)
	}

	render()
	go func() {
		defer close(stopped)
		for {
			select {
			case <-ticker.C:
				step++
				render()
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(done)
		<-stopped
	}
}

func (p *PatchExecutor) waitForRollback(ctx context.Context, ttl time.Duration) error {
	writer := p.output
	if writer == nil {
		writer = os.Stdout
	}

	rollbackAt := time.Now().Add(ttl)
	step := 0

	render := func() {
		remaining := time.Until(rollbackAt)
		countdown := formatCountdown(remaining)
		fmt.Fprintf(writer, "\r%s Waiting %s before rollback (at %s) - remaining %s", spinnerFrames[step%len(spinnerFrames)], ttl, rollbackAt.Format(time.RFC3339), countdown)
	}

	ticker := time.NewTicker(spinnerInterval)
	defer ticker.Stop()

	done := make(chan error, 1)
	go func() {
		done <- p.waitFunc(ctx, ttl)
	}()

	render()
	for {
		select {
		case err := <-done:
			fmt.Fprintln(writer)
			return err
		case <-ticker.C:
			step++
			render()
		}
	}
}

func formatCountdown(d time.Duration) string {
	seconds := int64(math.Ceil(d.Seconds()))
	if seconds < 0 {
		seconds = 0
	}
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60
	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, remainingSeconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, remainingSeconds)
}

func (p *PatchExecutor) printDryRunPreview(intent domain.PatchIntent, applyPayload map[string]any, rollbackPayload map[string]any) error {
	writer := p.output
	if writer == nil {
		writer = os.Stdout
	}

	fmt.Fprintln(writer, "Dry run enabled. No changes applied.")
	fmt.Fprintf(writer, "Target: %s/%s\n", intent.Resource.Kind, intent.Resource.Name)
	if intent.Resource.Namespace != "" {
		fmt.Fprintf(writer, "Namespace: %s\n", intent.Resource.Namespace)
	}
	if intent.Reason != "" {
		fmt.Fprintf(writer, "Reason: %s\n", intent.Reason)
	}
	if intent.TTL > 0 {
		fmt.Fprintf(writer, "TTL: %s\n", intent.TTL)
	}

	applyBytes, err := json.MarshalIndent(applyPayload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format apply patch: %w", err)
	}
	fmt.Fprintln(writer, "Apply patch:")
	fmt.Fprintln(writer, string(applyBytes))

	if intent.TTL <= 0 {
		fmt.Fprintln(writer, "Rollback patch: not applicable (ttl=0)")
		return nil
	}

	rollbackBytes, err := json.MarshalIndent(rollbackPayload, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format rollback patch: %w", err)
	}
	fmt.Fprintln(writer, "Rollback patch:")
	fmt.Fprintln(writer, string(rollbackBytes))
	return nil
}

func (p *PatchExecutor) resolveResource(kind string) (*meta.RESTMapping, error) {
	if err := p.ensureClients(); err != nil {
		return nil, err
	}

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
	if err := p.ensureClients(); err != nil {
		return nil, err
	}

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

type pathValue struct {
	path  []string
	value any
	found bool
}

func buildRollbackPayload(resource *unstructured.Unstructured, intent domain.PatchIntent) (map[string]any, error) {
	values := make([]pathValue, 0, len(intent.Patches)+2)

	for _, patch := range intent.Patches {
		value, found, err := unstructured.NestedFieldNoCopy(resource.Object, patch.Path...)
		if err != nil {
			return nil, fmt.Errorf("failed to read original value for %s: %w", strings.Join(patch.Path, "."), err)
		}
		if found {
			value = runtime.DeepCopyJSONValue(value)
		}

		values = append(values, pathValue{path: patch.Path, value: value, found: found})
	}

	if intent.Reason != "" {
		value, err := readAnnotationValue(resource, annotationReason)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	if intent.TTL != 0 {
		value, err := readAnnotationValue(resource, annotationTTL)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return buildMergePatchPayloadFromValues(values)
}

func readAnnotationValue(resource *unstructured.Unstructured, key string) (pathValue, error) {
	path := []string{"metadata", "annotations", key}
	value, found, err := unstructured.NestedFieldNoCopy(resource.Object, path...)
	if err != nil {
		return pathValue{}, fmt.Errorf("failed to read original annotation %q: %w", key, err)
	}
	if found {
		value = runtime.DeepCopyJSONValue(value)
	}
	return pathValue{path: path, value: value, found: found}, nil
}

func buildMergePatchPayloadFromValues(values []pathValue) (map[string]any, error) {
	root := map[string]any{}

	for _, item := range values {
		if len(item.path) == 0 {
			return nil, errors.New("patch path cannot be empty")
		}

		current := root

		for i, segment := range item.path {
			if i == len(item.path)-1 {
				if item.found {
					current[segment] = item.value
				} else {
					current[segment] = nil
				}
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
				return nil, fmt.Errorf("conflicting patch path at %s", strings.Join(item.path[:i+1], "."))
			}

			current = nested
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

func waitForTTL(ctx context.Context, ttl time.Duration) error {
	if ttl <= 0 {
		return nil
	}

	timer := time.NewTimer(ttl)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
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
