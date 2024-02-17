/*
Copyright 2024 Jason Kölker.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// +kubebuilder:object:generate=true
// +groupName=network.natpmp.jkoelker.github.io
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "network.natpmp.jkoelker.github.io"
	VersionV1 = "v1"
	Kind      = "NatPMP"
)

// GroupVersion returns the GroupVersion for the natpmp API.
func GroupVersion() schema.GroupVersion {
	return schema.GroupVersion{Group: GroupName, Version: VersionV1}
}

// GroupVersionKind returns the GroupVersionKind for the natpmp API.
func GroupVersionKind() schema.GroupVersionKind {
	return GroupVersion().WithKind(Kind)
}

// Adds the list of known types to the given scheme.
func AddToScheme(scheme *runtime.Scheme) error {
	groupVersion := GroupVersion()

	scheme.AddKnownTypes(
		groupVersion,
		&NatPMP{},
		&NatPMPList{},
	)

	scheme.AddKnownTypes(
		groupVersion,
		&metav1.Status{},
	)

	metav1.AddToGroupVersion(scheme, groupVersion)

	return nil
}
