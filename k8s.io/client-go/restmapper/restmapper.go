package restmapper

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ShortcutExpander struct {
	meta.RESTMapper
}

func NewDeferredDiscoveryRESTMapper(_ interface{}) meta.RESTMapper {
	return meta.NewDefaultRESTMapper(nil)
}

func NewShortcutExpander(mapper meta.RESTMapper, _ interface{}) *ShortcutExpander {
	return &ShortcutExpander{RESTMapper: mapper}
}

func (s *ShortcutExpander) ResourceFor(res schema.GroupVersionResource) (schema.GroupVersionResource, error) {
	return s.RESTMapper.ResourceFor(res)
}

func (s *ShortcutExpander) KindFor(res schema.GroupVersionResource) (schema.GroupVersionKind, error) {
	return s.RESTMapper.KindFor(res)
}

func (s *ShortcutExpander) RESTMapping(gk schema.GroupKind, versions ...string) (*meta.RESTMapping, error) {
	return s.RESTMapper.RESTMapping(gk, versions...)
}
