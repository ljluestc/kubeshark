package kubernetes

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateRESTMapper creates a REST mapper with shortcut support
func CreateRESTMapper(config *rest.Config) (meta.RESTMapper, error) {
	// Create the discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create discovery client: %w", err)
	}

	// Create a cached discovery client
	cachedDiscoveryClient, err := disk.NewCachedDiscoveryClientForConfig(
		config,
		clientcmd.RecommendedHomeDir+"/.kube/cache",
		clientcmd.RecommendedHomeDir+"/.kube/cache",
		180, // default expiry
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create cached discovery client: %w", err)
	}

	// Create the mapper
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cachedDiscoveryClient)

	// Use our custom shortcut expander helper function
	return NewShortcutExpander(mapper, discoveryClient)
}
