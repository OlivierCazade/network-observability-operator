/*
Copyright The Kubernetes Authors.

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

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v2beta2

import (
	v2beta2 "k8s.io/api/autoscaling/v2beta2"
)

// MetricSpecApplyConfiguration represents a declarative configuration of the MetricSpec type for use
// with apply.
type MetricSpecApplyConfiguration struct {
	Type              *v2beta2.MetricSourceType                        `json:"type,omitempty"`
	Object            *ObjectMetricSourceApplyConfiguration            `json:"object,omitempty"`
	Pods              *PodsMetricSourceApplyConfiguration              `json:"pods,omitempty"`
	Resource          *ResourceMetricSourceApplyConfiguration          `json:"resource,omitempty"`
	ContainerResource *ContainerResourceMetricSourceApplyConfiguration `json:"containerResource,omitempty"`
	External          *ExternalMetricSourceApplyConfiguration          `json:"external,omitempty"`
}

// MetricSpecApplyConfiguration constructs a declarative configuration of the MetricSpec type for use with
// apply.
func MetricSpec() *MetricSpecApplyConfiguration {
	return &MetricSpecApplyConfiguration{}
}

// WithType sets the Type field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Type field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithType(value v2beta2.MetricSourceType) *MetricSpecApplyConfiguration {
	b.Type = &value
	return b
}

// WithObject sets the Object field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Object field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithObject(value *ObjectMetricSourceApplyConfiguration) *MetricSpecApplyConfiguration {
	b.Object = value
	return b
}

// WithPods sets the Pods field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Pods field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithPods(value *PodsMetricSourceApplyConfiguration) *MetricSpecApplyConfiguration {
	b.Pods = value
	return b
}

// WithResource sets the Resource field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Resource field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithResource(value *ResourceMetricSourceApplyConfiguration) *MetricSpecApplyConfiguration {
	b.Resource = value
	return b
}

// WithContainerResource sets the ContainerResource field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ContainerResource field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithContainerResource(value *ContainerResourceMetricSourceApplyConfiguration) *MetricSpecApplyConfiguration {
	b.ContainerResource = value
	return b
}

// WithExternal sets the External field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the External field is set to the value of the last call.
func (b *MetricSpecApplyConfiguration) WithExternal(value *ExternalMetricSourceApplyConfiguration) *MetricSpecApplyConfiguration {
	b.External = value
	return b
}
