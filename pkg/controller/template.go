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
	"bytes"
	"errors"
	"fmt"
	"io"
	texttemplate "text/template"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	networkv1 "github.com/jkoelker/natpmp-controller/api/v1"
)

const decodeBufferSize = 4096

// TemplateSpec is a template safe version of the NatPMP spec object.
type TemplateSpec struct {
	ExternalPort int
	InternalPort int
	Lifetime     int
	Gateway      string
	Protocol     string
}

// TemplateStatus is a template safe version of the NatPMP status object.
type TemplateStatus struct {
	ExternalIP               string
	MappedInternalPort       int
	MappedExternalPort       int
	MappedLifetime           int
	SecondsSinceStartOfEpoch int
}

// TemplateDot is a template safe version of the NatPMP object.
type TemplateDot struct {
	Spec   TemplateSpec
	Status TemplateStatus
}

// ProcessTemplate takes a template string and a NatPMP object and returns
// a list of unstructured objects. The template is either a yaml or json
// template.
func ProcessTemplate(
	template string,
	natpmpCR networkv1.NatPMP,
) ([]*unstructured.Unstructured, error) {
	engine, err := texttemplate.New("natpmp").Parse(template)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	dot := TemplateDot{
		Spec: TemplateSpec{
			ExternalPort: natpmpCR.Spec.ExternalPort,
			InternalPort: natpmpCR.Spec.InternalPort,
			Lifetime:     natpmpCR.Spec.Lifetime,
			Gateway:      natpmpCR.Spec.Gateway,
			Protocol:     natpmpCR.Spec.Protocol,
		},
		Status: TemplateStatus{
			ExternalIP:               natpmpCR.Status.ExternalIP,
			MappedInternalPort:       natpmpCR.Status.MappedInternalPort,
			MappedExternalPort:       natpmpCR.Status.MappedExternalPort,
			MappedLifetime:           natpmpCR.Status.MappedLifetime,
			SecondsSinceStartOfEpoch: natpmpCR.Status.SecondsSinceStartOfEpoch,
		},
	}

	var buf bytes.Buffer
	if err := engine.Execute(&buf, dot); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	dec := yaml.NewYAMLOrJSONDecoder(&buf, decodeBufferSize)

	var objects []*unstructured.Unstructured

	for {
		var object unstructured.Unstructured

		if err := dec.Decode(&object); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to decode object: %w", err)
		}

		objects = append(objects, &object)
	}

	return objects, nil
}

// ProcessTemplates processes all templates for a NatPMP object.
func ProcessTemplates(natpmpCR networkv1.NatPMP) ([]*unstructured.Unstructured, error) {
	var objects []*unstructured.Unstructured

	for _, template := range natpmpCR.Spec.Templates {
		templateObjects, err := ProcessTemplate(template, natpmpCR)
		if err != nil {
			return nil, fmt.Errorf("failed to process template: %w", err)
		}

		objects = append(objects, templateObjects...)
	}

	return objects, nil
}
