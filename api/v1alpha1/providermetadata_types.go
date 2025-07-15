/*
Copyright 2024.

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

type Icon struct {
	Light Image `json:"light"`
	Dark  Image `json:"dark"`
}

type URL struct {
	URL string `json:"url,omitempty"`
}

type Image struct {
	URL  string `json:"url,omitempty"`
	Data string `json:"data,omitempty"`
}

type Link struct {
	DisplayName string `json:"displayName,omitempty"`
	URL         string `json:"URL,omitempty"`
}

type Contact struct {
	DisplayName string   `json:"displayName"`
	Email       string   `json:"email,omitempty"`
	Role        []string `json:"role,omitempty"`
}

// ProviderMetadataSpec defines the desired state of ProviderMetadata.
type ProviderMetadataSpec struct {
	Tags              []string `json:"tags,omitempty"`
	InstallableScopes []string `json:"installableScopes"`

	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`

	ConfigurationMetadata string    `json:"configurationMetadata,omitempty"`
	Contacts              []Contact `json:"contacts,omitempty"`
	Documentation         *URL      `json:"documentation,omitempty"`
	Icon                  *Icon     `json:"icon,omitempty"`

	Links                    []Link `json:"links,omitempty"`
	PreferredSupportChannels []Link `json:"preferredSupportChannels,omitempty"`
	HelpCenterData           []Link `json:"helpCenterData,omitempty"`
}

// ProviderMetadataStatus defines the observed state of ProviderMetadata.
type ProviderMetadataStatus struct{}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:resource:path=providermetadatas

// ProviderMetadata is the Schema for the providermetadata API.
type ProviderMetadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderMetadataSpec   `json:"spec,omitempty"`
	Status ProviderMetadataStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProviderMetadataList contains a list of ProviderMetadata.
type ProviderMetadataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderMetadata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderMetadata{}, &ProviderMetadataList{})
}
