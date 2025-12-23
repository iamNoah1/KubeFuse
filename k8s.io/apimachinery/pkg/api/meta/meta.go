package meta

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type RESTScopeName string

const (
	RESTScopeNameNamespace RESTScopeName = "namespace"
	RESTScopeNameRoot      RESTScopeName = "root"
)

type RESTScope interface {
	Name() RESTScopeName
}

type RESTMapping struct {
	Resource         schema.GroupVersionResource
	GroupVersionKind schema.GroupVersionKind
	Scope            RESTScope
}

type RESTMapper interface {
	ResourceFor(schema.GroupVersionResource) (schema.GroupVersionResource, error)
	KindFor(schema.GroupVersionResource) (schema.GroupVersionKind, error)
	RESTMapping(schema.GroupKind, ...string) (*RESTMapping, error)
}

type DefaultRESTMapper struct {
	mappings map[string]*RESTMapping
}

type restScopeNamespace struct{}

func (restScopeNamespace) Name() RESTScopeName { return RESTScopeNameNamespace }

type restScopeRoot struct{}

func (restScopeRoot) Name() RESTScopeName { return RESTScopeNameRoot }

var (
	RESTScopeNamespace RESTScope = restScopeNamespace{}
	RESTScopeRoot      RESTScope = restScopeRoot{}
)

func NewDefaultRESTMapper(_ []schema.GroupVersion) *DefaultRESTMapper {
	return &DefaultRESTMapper{mappings: map[string]*RESTMapping{}}
}

func (m *DefaultRESTMapper) AddSpecific(gvk schema.GroupVersionKind, gvr schema.GroupVersionResource, _ schema.GroupVersionResource, scope RESTScope) {
	key := gvrKey(gvr)
	m.mappings[key] = &RESTMapping{Resource: gvr, GroupVersionKind: gvk, Scope: scope}
}

func gvrKey(gvr schema.GroupVersionResource) string {
	return gvr.Group + "/" + gvr.Version + "/" + gvr.Resource
}

func (m *DefaultRESTMapper) ResourceFor(gvr schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	if mapping, ok := m.mappings[gvrKey(gvr)]; ok {
		return mapping.Resource, nil
	}
	for _, mapping := range m.mappings {
		if mapping.Resource.Resource == gvr.Resource {
			return mapping.Resource, nil
		}
	}
	return gvr, nil
}

func (m *DefaultRESTMapper) KindFor(gvr schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	if mapping, ok := m.mappings[gvrKey(gvr)]; ok {
		return mapping.GroupVersionKind, nil
	}
	for _, mapping := range m.mappings {
		if mapping.Resource.Resource == gvr.Resource {
			return mapping.GroupVersionKind, nil
		}
	}
	return schema.GroupVersionKind{Group: gvr.Group, Version: gvr.Version, Kind: gvr.Resource}, nil
}

func (m *DefaultRESTMapper) RESTMapping(gk schema.GroupKind, versions ...string) (*RESTMapping, error) {
	for _, mapping := range m.mappings {
		if mapping.GroupVersionKind.Group == gk.Group && mapping.GroupVersionKind.Kind == gk.Kind {
			if len(versions) == 0 || mapping.GroupVersionKind.Version == versions[0] {
				return mapping, nil
			}
		}
	}
	for _, mapping := range m.mappings {
		if strings.EqualFold(mapping.GroupVersionKind.Kind, gk.Kind) || strings.EqualFold(mapping.Resource.Resource, gk.Kind) {
			return mapping, nil
		}
	}
	return nil, &NoKindMatchError{GroupKind: gk}
}

type NoKindMatchError struct {
	GroupKind schema.GroupKind
}

func (e *NoKindMatchError) Error() string { return "no matches for kind " + e.GroupKind.Kind }
