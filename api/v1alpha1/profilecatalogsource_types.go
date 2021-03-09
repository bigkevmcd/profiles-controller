/*
Copyright 2021.

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

// Profile might be able to be reused from the Profile specification.
type Profile struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
	Maturity    string `json:"maturity,omitempty"` // This should be annotated to reflect correct values.
	Publisher   string `json:"publisher,omitempty"`
}

// ProfileCatalogSourceSpec defines the desired state of ProfileCatalogSource.
type ProfileCatalogSourceSpec struct {
	// These are purely metadata fields.
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
	Publisher   string `json:"publisher,omitempty"`

	EmbeddedProfiles []Profile `json:"profiles,omitempty"`
}

// ProfileCatalogSourceStatus defines the observed state of
// ProfileCatalogSource.
type ProfileCatalogSourceStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProfileCatalogSource is the Schema for the profilecatalogsources API.
type ProfileCatalogSource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileCatalogSourceSpec   `json:"spec,omitempty"`
	Status ProfileCatalogSourceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileCatalogSourceList contains a list of ProfileCatalogSource
type ProfileCatalogSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProfileCatalogSource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProfileCatalogSource{}, &ProfileCatalogSourceList{})
}
