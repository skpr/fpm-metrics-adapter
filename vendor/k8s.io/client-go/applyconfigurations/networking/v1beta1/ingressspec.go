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

package v1beta1

// IngressSpecApplyConfiguration represents an declarative configuration of the IngressSpec type for use
// with apply.
type IngressSpecApplyConfiguration struct {
	IngressClassName *string                           `json:"ingressClassName,omitempty"`
	Backend          *IngressBackendApplyConfiguration `json:"backend,omitempty"`
	TLS              []IngressTLSApplyConfiguration    `json:"tls,omitempty"`
	Rules            []IngressRuleApplyConfiguration   `json:"rules,omitempty"`
}

// IngressSpecApplyConfiguration constructs an declarative configuration of the IngressSpec type for use with
// apply.
func IngressSpec() *IngressSpecApplyConfiguration {
	return &IngressSpecApplyConfiguration{}
}

// WithIngressClassName sets the IngressClassName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the IngressClassName field is set to the value of the last call.
func (b *IngressSpecApplyConfiguration) WithIngressClassName(value string) *IngressSpecApplyConfiguration {
	b.IngressClassName = &value
	return b
}

// WithBackend sets the Backend field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Backend field is set to the value of the last call.
func (b *IngressSpecApplyConfiguration) WithBackend(value *IngressBackendApplyConfiguration) *IngressSpecApplyConfiguration {
	b.Backend = value
	return b
}

// WithTLS adds the given value to the TLS field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the TLS field.
func (b *IngressSpecApplyConfiguration) WithTLS(values ...*IngressTLSApplyConfiguration) *IngressSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTLS")
		}
		b.TLS = append(b.TLS, *values[i])
	}
	return b
}

// WithRules adds the given value to the Rules field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Rules field.
func (b *IngressSpecApplyConfiguration) WithRules(values ...*IngressRuleApplyConfiguration) *IngressSpecApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithRules")
		}
		b.Rules = append(b.Rules, *values[i])
	}
	return b
}
