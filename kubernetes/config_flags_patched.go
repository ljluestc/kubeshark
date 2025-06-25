package kubernetes

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
)

// PatchedConfigFlags is a wrapper around ConfigFlags that provides compatibility
// with different Kubernetes API versions
type PatchedConfigFlags struct {
	*genericclioptions.ConfigFlags
}

// ToRESTMapper returns a mapper with compatibility for different K8s versions
func (f *PatchedConfigFlags) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	// Use our custom shortcut expander that handles different API versions
	return NewShortcutExpander(mapper, discoveryClient), nil
}

// NewPatchedConfigFlags creates ConfigFlags with version compatibility fixes
func NewPatchedConfigFlags(usePersistentConfig bool) *PatchedConfigFlags {
	return &PatchedConfigFlags{
		ConfigFlags: genericclioptions.NewConfigFlags(usePersistentConfig),
	}
}
