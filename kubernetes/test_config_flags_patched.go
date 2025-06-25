package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
)

// PatchedTestConfigFlags is a wrapper around TestConfigFlags that provides compatibility
// with different Kubernetes API versions
type PatchedTestConfigFlags struct {
	*genericclioptions.TestConfigFlags
}

// ToRESTMapper implements RESTClientGetter with version compatibility
func (f *PatchedTestConfigFlags) ToRESTMapper() (meta.RESTMapper, error) {
	if f.RESTMapper != nil {
		return f.RESTMapper, nil
	}
	if f.DiscoveryClient != nil {
		mapper := restmapper.NewDeferredDiscoveryRESTMapper(f.DiscoveryClient)
		// Use our custom NewShortcutExpander function that includes the required third argument
		return NewShortcutExpander(mapper, f.DiscoveryClient), nil
	}
	return nil, fmt.Errorf("no restmapper")
}

// NewPatchedTestConfigFlags creates a TestConfigFlags with version compatibility fixes
func NewPatchedTestConfigFlags() *PatchedTestConfigFlags {
	return &PatchedTestConfigFlags{
		TestConfigFlags: genericclioptions.NewTestConfigFlags(),
	}
}
