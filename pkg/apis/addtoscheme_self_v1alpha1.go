package apis

import (
	api "awsctrl.io/pkg/apis/self/v1alpha1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, api.SchemeBuilder.AddToScheme)
}
