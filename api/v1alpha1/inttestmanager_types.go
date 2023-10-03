/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IntTestManagerSpec defines the desired state of IntTestManager
type IntTestManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Pod template to inject into the DM Worker Daemonset to allow for
	// verification of data movement through the direct use of the volumes
	// attached to the DM worker pods.
	// DMHelperTemplate corev1.PodTemplateSpec `json:"template"`
}

// IntTestManagerStatus defines the observed state of IntTestManager
type IntTestManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IntTestManager is the Schema for the inttestmanagers API
type IntTestManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IntTestManagerSpec   `json:"spec,omitempty"`
	Status IntTestManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IntTestManagerList contains a list of IntTestManager
type IntTestManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IntTestManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IntTestManager{}, &IntTestManagerList{})
}
