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

package controller

import (
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	networkv1 "github.com/jkoelker/natpmp-controller/api/v1"
)

const (
	TCP = "tcp"
	UDP = "udp"
)

// IsValidProtocol returns true if the protocol is valid. Valid protocols are
// TCP and UDP.
func IsValidProtocol(protocol string) bool {
	return protocol == TCP || protocol == UDP
}

// ValidateGateway returns an error if the gateway is not a valid IP address.
func ValidateGateway(gateway string, path ...string) (net.IP, *field.Error) {
	ip := net.ParseIP(gateway)
	if ip == nil {
		return nil, field.Invalid(field.NewPath("spec", path...), gateway, "invalid IP address")
	}

	return ip, nil
}

// ValidateLifetime returns an error if the lifetime is less than 1.
func ValidateLifetime(lifetime int, path ...string) *field.Error {
	if lifetime < 1 {
		return field.Invalid(field.NewPath("spec", path...), lifetime, "invalid lifetime")
	}

	return nil
}

// ValidatePort returns an error if the port is not between 0 and 65535.
func ValidatePort(port int, path ...string) *field.Error {
	if port < 0 || port > 65535 {
		return field.Invalid(field.NewPath("spec", path...), port, "invalid port")
	}

	return nil
}

// ValidateProtocol returns an error if the protocol is not TCP or UDP.
func ValidateProtocol(protocol string, path ...string) *field.Error {
	if !IsValidProtocol(protocol) {
		return field.Invalid(field.NewPath("spec", path...), protocol, "invalid protocol")
	}

	return nil
}

// ValidateNatPMP returns the gateway IP, protocol, and a list of errors if any.
func ValidateNatPMP(natpmpCR networkv1.NatPMP) (net.IP, string, field.ErrorList) {
	var allErrs field.ErrorList

	gateway, err := ValidateGateway(natpmpCR.Spec.Gateway, "gateway")
	if err != nil {
		allErrs = append(allErrs, err)
	}

	protocol := strings.ToLower(natpmpCR.Spec.Protocol)
	if err := ValidateProtocol(protocol, "protocol"); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := ValidatePort(natpmpCR.Spec.ExternalPort, "externalPort"); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := ValidatePort(natpmpCR.Spec.InternalPort, "internalPort"); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := ValidateLifetime(natpmpCR.Spec.Lifetime, "lifetime"); err != nil {
		allErrs = append(allErrs, err)
	}

	return gateway, protocol, allErrs
}
