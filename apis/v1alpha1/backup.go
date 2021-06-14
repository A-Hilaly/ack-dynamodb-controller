// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BackupSpec defines the desired state of Backup
type BackupSpec struct {
	// Specified name for the backup.
	// +kubebuilder:validation:Required
	BackupName *string `json:"backupName"`
	// The name of the table.
	// +kubebuilder:validation:Required
	TableName *string `json:"tableName"`
}

// BackupStatus defines the observed state of Backup
type BackupStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// Time at which the backup was created. This is the request time of the backup.
	BackupCreationDateTime *metav1.Time `json:"backupCreationDateTime,omitempty"`
	// Time at which the automatic on-demand backup created by DynamoDB will expire.
	// This SYSTEM on-demand backup expires automatically 35 days after its creation.
	BackupExpiryDateTime *metav1.Time `json:"backupExpiryDateTime,omitempty"`
	// Size of the backup in bytes.
	BackupSizeBytes *int64 `json:"backupSizeBytes,omitempty"`
	// Backup can be in one of the following states: CREATING, ACTIVE, DELETED.
	BackupStatus *string `json:"backupStatus,omitempty"`
	// BackupType:
	//
	//    * USER - You create and manage these using the on-demand backup feature.
	//
	//    * SYSTEM - If you delete a table with point-in-time recovery enabled,
	//    a SYSTEM backup is automatically created and is retained for 35 days (at
	//    no additional cost). System backups allow you to restore the deleted table
	//    to the state it was in just before the point of deletion.
	//
	//    * AWS_BACKUP - On-demand backup created by you from AWS Backup service.
	BackupType *string `json:"backupType,omitempty"`
}

// Backup is the Schema for the Backups API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Backup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BackupSpec   `json:"spec,omitempty"`
	Status            BackupStatus `json:"status,omitempty"`
}

// BackupList contains a list of Backup
// +kubebuilder:object:root=true
type BackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backup{}, &BackupList{})
}