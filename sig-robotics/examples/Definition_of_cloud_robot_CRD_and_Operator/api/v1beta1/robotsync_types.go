/*
Copyright 2022 The KubeEdge Authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RobotSyncSpec defines the desired state of RobotSync
type RobotSyncSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

}

// RobotSyncStatus defines the observed state of RobotSync
type RobotSyncStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastHeartbeat  map[string]metav1.Time `json:"lastHeartbeat,omitempty"` // RobotID is contained in heatbeat packet, which is used to correspondence between Robot name and RobotID
	RegistedRobots []uint                 `json:"registedRobots,omitempty"`	// List of registered robots
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RobotSync is the Schema for the robotsyncs API
type RobotSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RobotSyncSpec   `json:"spec,omitempty"`
	Status RobotSyncStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RobotSyncList contains a list of RobotSync
type RobotSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RobotSync `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RobotSync{}, &RobotSyncList{})
}
