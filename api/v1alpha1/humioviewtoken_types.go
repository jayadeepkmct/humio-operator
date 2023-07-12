/*
Copyright 2020 Humio https://humio.com

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

type ViewTokenPermissions struct {
	ChangeUserAccess                  bool `json:"changeUserAccess,omitempty"`
	ChangeTriggersAndActions          bool `json:"changeTriggersAndActions,omitempty"`
	ChangeDashboards                  bool `json:"changeDashboards,omitempty"`
	ChangeDashboardReadonlyToken      bool `json:"changeDashboardReadonlyToken,omitempty"`
	ChangeFiles                       bool `json:"changeFiles,omitempty"`
	ChangeInteractions                bool `json:"changeInteractions,omitempty"`
	ChangeParsers                     bool `json:"changeParsers,omitempty"`
	ChangeSavedQueries                bool `json:"changeSavedQueries,omitempty"`
	ConnectView                       bool `json:"connectView,omitempty"`
	ChangeDataDeletionPermissions     bool `json:"changeDataDeletionPermissions,omitempty"`
	ChangeRetention                   bool `json:"changeRetention,omitempty"`
	ChangeDefaultSearchSettings       bool `json:"changeDefaultSearchSettings,omitempty"`
	ChangeS3ArchivingSettings         bool `json:"changeS3ArchivingSettings,omitempty"`
	DeleteDataSources                 bool `json:"deleteDataSources,omitempty"`
	DeleteRepositoryOrView            bool `json:"deleteRepositoryOrView,omitempty"`
	DeleteEvents                      bool `json:"deleteEvents,omitempty"`
	ReadAccess                        bool `json:"readAccess,omitempty"`
	ChangeIngestTokens                bool `json:"changeIngestTokens,omitempty"`
	ChangePackages                    bool `json:"changePackages,omitempty"`
	ChangeViewOrRepositoryDescription bool `json:"changeViewOrRepositoryDescription,omitempty"`
	ChangeConnections                 bool `json:"changeConnections,omitempty"`
	EventForwarding                   bool `json:"eventForwarding,omitempty"`
	QueryDashboard                    bool `json:"queryDashboard,omitempty"`
	ChangeViewOrRepositoryPermissions bool `json:"changeViewOrRepositoryPermissions,omitempty"`
	ChangeFdrFeeds                    bool `json:"changeFdrFeeds,omitempty"`
	OrganizationOwnedQueries          bool `json:"organizationOwnedQueries,omitempty"`
	ReadExternalFunctions             bool `json:"readExternalFunctions,omitempty"`
}

// HumioViewTokenSpec defines the desired state of HumioViewToken
type HumioViewTokenSpec struct {
	// ManagedClusterName refers to an object of type HumioCluster that is managed by the operator where the Humio
	// resources should be created.
	// This conflicts with ExternalClusterName.
	ManagedClusterName string `json:"managedClusterName,omitempty"`
	// ExternalClusterName refers to an object of type HumioExternalCluster where the Humio resources should be created.
	// This conflicts with ManagedClusterName.
	ExternalClusterName string `json:"externalClusterName,omitempty"`
	// Name is the name of the ingest token inside Humio
	Name string `json:"name"`
	// ViewName is the name of the Humio view under which the token should be created
	ViewName string `json:"viewName,omitempty"`
	// TokenSecretName specifies the name of the Kubernetes secret that will be created
	// and contain the ingest token. The key in the secret storing the ingest token is "token".
	// This field is optional.
	TokenSecretName string `json:"tokenSecretName,omitempty"`
	// TokenSecretLabels specifies additional key,value pairs to add as labels on the Kubernetes Secret containing
	// the ingest token.
	// This field is optional.
	TokenSecretLabels map[string]string `json:"tokenSecretLabels,omitempty"`
	// Permissions indicates what all permissions to be granted for this token
	Permissions ViewTokenPermissions `json:"permissions,omitempty"`
}

// HumioViewTokenStatus defines the observed state of HumioViewToken
type HumioViewTokenStatus struct {
	// State reflects the current state of the HumioViewToken
	State string `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=humioviewtokens,scope=Namespaced
//+kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.state",description="The state of the view token"
//+operator-sdk:gen-csv:customresourcedefinitions.displayName="Humio View Token"

// HumioViewToken is the Schema for the humioviewtokens API
type HumioViewToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HumioViewTokenSpec   `json:"spec,omitempty"`
	Status HumioViewTokenStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HumioViewTokenList contains a list of HumioViewToken
type HumioViewTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HumioViewToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HumioViewToken{}, &HumioViewTokenList{})
}
