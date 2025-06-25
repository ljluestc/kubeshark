package kubernetes

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/restmapper"
)

// NewShortcutExpander creates a shortcut expander with the appropriate arguments
// This handles the API change in NewShortcutExpander between Kubernetes versions
func NewShortcutExpander(mapper meta.RESTMapper, discoveryClient discovery.DiscoveryInterface) meta.RESTMapper {
	// For k8s.io/cli-runtime v0.26.x, we use the 2-argument version
	return restmapper.NewShortcutExpander(mapper, discoveryClient)
}

```
}
}
