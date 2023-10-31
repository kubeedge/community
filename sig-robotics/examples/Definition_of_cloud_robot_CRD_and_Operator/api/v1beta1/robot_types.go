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

// RobotSpec defines the desired state of Robot
type RobotSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	RobotID uint `json:"robotID,omitempty"`	// Globally unique identifier of robot
}

// RobotStatus defines the observed state of Robot
type RobotStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Sensors        []Sensor         `json:"sensors,omitempty"`	// List of sensors equipped on the robot
	BatteryStatus  BatteryStatus    `json:"batteryStatus,omitempty"`
	ResourceStatus ResourceStatus   `json:"resourceStatus,omitempty"`
	Position       Position         `json:"position,omitempty"`
	RunningStatus  RunningStatus    `json:"runningStatus,omitempty"`
	UnderTask      bool             `json:"underTask,omitempty"`	// If robot is under a task
	TaskInfo       TaskInfo         `json:"taskInfo,omitempty"`	// Include the detail task information
	AbnormalEvents []AbnormalEvents `json:"abnormalEvents,omitempty"`	// Exception information of robot
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Robot is the Schema for the robots API
type Robot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RobotSpec   `json:"spec,omitempty"`
	Status RobotStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RobotList contains a list of Robot
type RobotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Robot `json:"items"`
}

type Position struct {
	SegmentStartPointID *uint    `json:"segmentStartPointID,omitempty"` //  start point ID of segment
	SegmentEndPointID   *uint    `json:"segmentEndPointID,omitempty"`
	X                   *float64 `json:"x,omitempty"` // mm
	Y                   *float64 `json:"y,omitempty"` // mm
	PreviousPoint       *uint    `json:"previousPoint,omitempty"`
}

type RunningStatus struct {
	LinearVelocity  *int     `json:"linearVelocity,omitempty"`  // mm/s
	AngularVelocity *float32 `json:"angularVelocity,omitempty"` // rad/s
	IMRStatus       *uint16  `json:"imrStatus,omitempty"`       // 0x00: available, 0x01: running, 0x02: stopped, 0x03: physically offline, 0x04: logically offline
}

type TaskInfo struct {
	OrderID              *uint                  `json:"orderID,omitempty"`
	TaskID               *uint                  `json:"taskID,omitempty"`
	PointStateSequence   []PointStateSequence   `json:"pointStateSequence,omitempty"` // received but not executed points
	SegmentStateSequence []SegmentStateSequence `json:"segmentStateSequence,omitempty"`
	RequiredSensors      []Sensor               `json:"requiredSensors,omitempty"`
}

type BatteryStatus struct {
	BatterySerialNumber []byte  `json:"batterySerialNumber,omitempty"`
	BatteryStatus       uint16  `json:"batteryStatus,omitempty"` // 0x01: charging, 0x02: discharging
	PowerPercentage     float32 `json:"powerPercentage,omitempty"`
	ChargingTimes       uint    `json:"chargingTimes,omitempty"`
}

type AbnormalEvents struct {
	Identifier     uint16 `json:"identifier,omitempty"` // 0xFFFF
	EventCode      uint16 `json:"eventCode,omitempty"`
	ExceptionLevel uint16 `json:"exceptionLevel,omitempty"` // 0x00: info, 0x01: warning, 0x02: error
	Description    string `json:"description,omitempty"`
}

// first edition
type Sensor struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	// Data        interface{} `json:"data,omitempty"`
}

type ResourceStatus struct {
	CpuUsagePercentage float32 `json:"cpuUsagePercentage,omitempty"`
	GpuUsagePercentage float32 `json:"gpuUsagePercentage,omitempty"`
	MemUsage           float32 `json:"memUsage,omitempty"`         // MB
	MemAvailable       float32 `json:"memAvailable,omitempty"`     // MB
	DiskUsage          float32 `json:"diskUsage,omitempty"`        // MB
	DiskAvailable      float32 `json:"diskAvailable,omitempty"`    // MB
	NetworkBandwidth   float32 `json:"networkBandwidth,omitempty"` // MB/s
}

func init() {
	SchemeBuilder.Register(&Robot{}, &RobotList{})
}
