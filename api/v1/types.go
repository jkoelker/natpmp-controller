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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NatPMPSpec defines the desired state of NatPMP.
type NatPMPSpec struct {
	// ExternalPort is the requested external port number to map.
	ExternalPort int `json:"externalPort"`

	// InternalPort is the internal port number that the external port maps to.
	InternalPort int `json:"internalPort"`

	// Lifetime is the duration in seconds for which the port mapping should
	// be active.
	Lifetime int `json:"lifetime"`

	// Gateway is the address or identifier of the NAT-PMP gateway.
	Gateway string `json:"gateway"`

	// Protocol is the protocol for the port mapping (TCP/UDP).
	Protocol string `json:"protocol"`

	// Templates is the raw templates that will be used to create or update
	// resources via server-side apply. Each template must be a valid
	// Kubernetes YAML or JSON document. The templates will be applied in
	// order. The templates may reference the following variables:
	//
	//   * .Spec.ExternalPort
	//   * .Spec.InternalPort
	//   * .Spec.Protocol
	//   * .Spec.Gateway
	//   * .Spec.Lifetime
	//   * .Status.ExternalIP
	//   * .Status.MappedInternalPort
	//   * .Status.MappedExternalPort
	//   * .Status.MappedLifetime
	//   * .Status.SecondsSinceStartOfEpoch
	Templates []string `json:"templates"`
}

// NatPMPStatus defines the observed state of NatPMP.
type NatPMPStatus struct {
	// ExternalIP is the external IP address of the gateway.
	ExternalIP string `json:"externalIP,omitempty"`

	// MappedInternalPort is the internal port number that the external port maps to.
	MappedInternalPort int `json:"internalPort,omitempty"`

	// MappedExternalPort is the external port number that was successfully
	// mapped.
	MappedExternalPort int `json:"mappedExternalPort,omitempty"`

	// MappedLifetime is the duration in seconds for which the port mapping
	// will be active.
	MappedLifetime int `json:"mappedLifetime,omitempty"`

	// SecondsSinceStartOfEpoch is the number of seconds since the start of
	// the epoch.
	SecondsSinceStartOfEpoch int `json:"secondsSinceStartOfEpoch,omitempty"`

	// Conditions represent the latest available observations of the
	// resource's state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NatPMP is the Schema for the natpmps API.
type NatPMP struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NatPMPSpec   `json:"spec,omitempty"`
	Status NatPMPStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NatPMPList contains a list of NatPMP.
type NatPMPList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NatPMP `json:"items"`
}
