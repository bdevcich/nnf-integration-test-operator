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
	"github.com/DataWorkflowServices/dws/utils/updater"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IntTestHelperSpec defines the desired state of IntTestHelper
type IntTestHelperSpec struct {
	Command string `json:"command,omitempty"`
}

// IntTestHelperStatus defines the observed state of IntTestHelper
type IntTestHelperStatus struct {
	Command     string          `json:"command,omitempty"`
	Output      string          `json:"output,omitempty"`
	Error       string          `json:"error,omitempty"`
	ElapsedTime metav1.Duration `json:"elapsedTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IntTestHelper is the Schema for the inttesthelpers API
type IntTestHelper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IntTestHelperSpec   `json:"spec,omitempty"`
	Status IntTestHelperStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IntTestHelperList contains a list of IntTestHelper
type IntTestHelperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IntTestHelper `json:"items"`
}

func (i *IntTestHelper) GetStatus() updater.Status[*IntTestHelperStatus] {
	return &i.Status
}

func init() {
	SchemeBuilder.Register(&IntTestHelper{}, &IntTestHelperList{})
}
