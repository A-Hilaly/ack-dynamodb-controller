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

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TableSpec defines the desired state of Table.
type TableSpec struct {

	// An array of attributes that describe the key schema for the table and indexes.
	// +kubebuilder:validation:Required
	AttributeDefinitions   []*AttributeDefinition  `json:"attributeDefinitions"`
	BillingMode            *string                 `json:"billingMode,omitempty"`
	GlobalSecondaryIndexes []*GlobalSecondaryIndex `json:"globalSecondaryIndexes,omitempty"`

	//--
	KeySchema             []*KeySchemaElement    `json:"keySchema"`
	LocalSecondaryIndexes []*LocalSecondaryIndex `json:"localSecondaryIndexes,omitempty"`
	ProvisionedThroughput *ProvisionedThroughput `json:"provisionedThroughput,omitempty"`
	SSESpecification      *SSESpecification      `json:"sseSpecification,omitempty"`
	StreamSpecification   *StreamSpecification   `json:"streamSpecification,omitempty"`
	TableClass            *string                `json:"tableClass,omitempty"`
	TableName             *string                `json:"tableName"`
	// A list of key-value pairs to label the table. For more information, see Tagging
	// for DynamoDB (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Tagging.html).
	Tags []*Tag `json:"tags,omitempty"`
	// Represents the settings used to enable or disable Time to Live for the specified
	// table.
	TimeToLive *TimeToLiveSpecification `json:"timeToLive,omitempty"`
}

// TableStatus defines the observed state of Table
type TableStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// Contains information about the table archive.
	// +kubebuilder:validation:Optional
	ArchivalSummary *ArchivalSummary `json:"archivalSummary,omitempty"`
	// Contains the details for the read/write capacity mode.
	// +kubebuilder:validation:Optional
	BillingModeSummary *BillingModeSummary `json:"billingModeSummary,omitempty"`
	// The date and time when the table was created, in UNIX epoch time (http://www.epochconverter.com/)
	// format.
	// +kubebuilder:validation:Optional
	CreationDateTime *metav1.Time `json:"creationDateTime,omitempty"`
	// Represents the version of global tables (https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GlobalTables.html)
	// in use, if the table is replicated across Amazon Web Services Regions.
	// +kubebuilder:validation:Optional
	GlobalTableVersion *string `json:"globalTableVersion,omitempty"`
	// The number of items in the specified table. DynamoDB updates this value approximately
	// every six hours. Recent changes might not be reflected in this value.
	// +kubebuilder:validation:Optional
	ItemCount *int64 `json:"itemCount,omitempty"`
	// The Amazon Resource Name (ARN) that uniquely identifies the latest stream
	// for this table.
	// +kubebuilder:validation:Optional
	LatestStreamARN *string `json:"latestStreamARN,omitempty"`
	// A timestamp, in ISO 8601 format, for this stream.
	//
	// Note that LatestStreamLabel is not a unique identifier for the stream, because
	// it is possible that a stream from another table might have the same timestamp.
	// However, the combination of the following three elements is guaranteed to
	// be unique:
	//
	//    * Amazon Web Services customer ID
	//
	//    * Table name
	//
	//    * StreamLabel
	// +kubebuilder:validation:Optional
	LatestStreamLabel *string `json:"latestStreamLabel,omitempty"`
	// Represents replicas of the table.
	// +kubebuilder:validation:Optional
	Replicas []*ReplicaDescription `json:"replicas,omitempty"`
	// Contains details for the restore.
	// +kubebuilder:validation:Optional
	RestoreSummary *RestoreSummary `json:"restoreSummary,omitempty"`
	// The description of the server-side encryption status on the specified table.
	// +kubebuilder:validation:Optional
	SSEDescription *SSEDescription `json:"sseDescription,omitempty"`
	// Contains details of the table class.
	// +kubebuilder:validation:Optional
	TableClassSummary *TableClassSummary `json:"tableClassSummary,omitempty"`
	// Unique identifier for the table for which the backup was created.
	// +kubebuilder:validation:Optional
	TableID *string `json:"tableID,omitempty"`
	// The total size of the specified table, in bytes. DynamoDB updates this value
	// approximately every six hours. Recent changes might not be reflected in this
	// value.
	// +kubebuilder:validation:Optional
	TableSizeBytes *int64 `json:"tableSizeBytes,omitempty"`
	// The current state of the table:
	//
	//    * CREATING - The table is being created.
	//
	//    * UPDATING - The table is being updated.
	//
	//    * DELETING - The table is being deleted.
	//
	//    * ACTIVE - The table is ready for use.
	//
	//    * INACCESSIBLE_ENCRYPTION_CREDENTIALS - The KMS key used to encrypt the
	//    table in inaccessible. Table operations may fail due to failure to use
	//    the KMS key. DynamoDB will initiate the table archival process when a
	//    table's KMS key remains inaccessible for more than seven days.
	//
	//    * ARCHIVING - The table is being archived. Operations are not allowed
	//    until archival is complete.
	//
	//    * ARCHIVED - The table has been archived. See the ArchivalReason for more
	//    information.
	// +kubebuilder:validation:Optional
	TableStatus *string `json:"tableStatus,omitempty"`
}

// Table is the Schema for the Tables API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Table struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              TableSpec   `json:"spec,omitempty"`
	Status            TableStatus `json:"status,omitempty"`
}

// TableList contains a list of Table
// +kubebuilder:object:root=true
type TableList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Table `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Table{}, &TableList{})
}
