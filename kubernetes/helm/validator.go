package helm
package helm

import (
	"fmt"
)

// ValidateHelmRelease checks if a Helm release is valid
func ValidateHelmRelease(releaseName string, namespace string) error {
	if releaseName == "" {
		return fmt.Errorf("release name cannot be empty")
	}
	
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	
	// This is a placeholder for actual validation logic
	// In a real implementation, we would check if the release exists in the cluster
	
	return nil
}

// ValidateHelmChart checks if a Helm chart is valid
func ValidateHelmChart(chartName string, version string) error {
	if chartName == "" {
		return fmt.Errorf("chart name cannot be empty")
	}
	
	// This is a placeholder for actual validation logic
	// In a real implementation, we would check if the chart exists in the repository
	
	return nil
}
import (
	"fmt"
)

// ValidateHelmRelease checks if a Helm release is valid
func ValidateHelmRelease(releaseName string, namespace string) error {
	if releaseName == "" {
		return fmt.Errorf("release name cannot be empty")
	}
	
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	
	// This is a placeholder for actual validation logic
	// In a real implementation, we would check if the release exists in the cluster
	
	return nil
}

// ValidateHelmChart checks if a Helm chart is valid
func ValidateHelmChart(chartName string, version string) error {
	if chartName == "" {
		return fmt.Errorf("chart name cannot be empty")
	}
	
	// This is a placeholder for actual validation logic
	// In a real implementation, we would check if the chart exists in the repository
	
	return nil
}
