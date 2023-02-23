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

package table

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/dynamodb"
	corev1 "k8s.io/api/core/v1"

	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
)

var (
	ErrTableDeleting = errors.New("Table in 'DELETING' state, cannot be modified or deleted")
	ErrTableCreating = errors.New("Table in 'CREATING' state, cannot be modified or deleted")
	ErrTableUpdating = errors.New("Table in 'UPDATING' state, cannot be modified or deleted")
)

var (
	// TerminalStatuses are the status strings that are terminal states for a
	// DynamoDB table
	TerminalStatuses = []v1alpha1.TableStatus_SDK{
		v1alpha1.TableStatus_SDK_ARCHIVING,
		v1alpha1.TableStatus_SDK_DELETING,
	}
)

var DefaultTTLEnabledValue = false

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		ErrTableDeleting,
		5*time.Second,
	)
	requeueWaitWhileCreating = ackrequeue.NeededAfter(
		ErrTableCreating,
		5*time.Second,
	)
	requeueWaitWhileUpdating = ackrequeue.NeededAfter(
		ErrTableUpdating,
		5*time.Second,
	)
)

// tableHasTerminalStatus returns whether the supplied Dynamodb table is in a
// terminal state
func tableHasTerminalStatus(r *resource) bool {
	if r.ko.Status.TableStatus == nil {
		return false
	}
	ts := *r.ko.Status.TableStatus
	for _, s := range TerminalStatuses {
		if ts == string(s) {
			return true
		}
	}
	return false
}

// isTableCreating returns true if the supplied DynamodbDB table is in the process
// of being created
func isTableCreating(r *resource) bool {
	if r.ko.Status.TableStatus == nil {
		return false
	}
	dbis := *r.ko.Status.TableStatus
	return dbis == string(v1alpha1.TableStatus_SDK_CREATING)
}

// isTableDeleting returns true if the supplied DynamodbDB table is in the process
// of being deleted
func isTableDeleting(r *resource) bool {
	if r.ko.Status.TableStatus == nil {
		return false
	}
	dbis := *r.ko.Status.TableStatus
	return dbis == string(v1alpha1.TableStatus_SDK_DELETING)
}

// isTableUpdating returns true if the supplied DynamodbDB table is in the process
// of being deleted
func isTableUpdating(r *resource) bool {
	if r.ko.Status.TableStatus == nil {
		return false
	}
	dbis := *r.ko.Status.TableStatus
	return dbis == string(v1alpha1.TableStatus_SDK_UPDATING)
}

func (rm *resourceManager) customUpdateTable(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	fmt.Println("STARTING SSE ======================")
	x, _ := json.Marshal(desired.ko.Spec.SSESpecification)
	fmt.Println("SSE HOOK A", string(x))
	x, _ = json.Marshal(latest.ko.Spec.SSESpecification)
	fmt.Println("SSE HOOK B", string(x))

	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdateTable")
	defer func(err error) { exit(err) }(err)

	if isTableDeleting(latest) {
		msg := "table is currently being deleted"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhileDeleting
	}
	if isTableCreating(latest) {
		msg := "table is currently being created"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhileCreating
	}
	if isTableUpdating(latest) {
		msg := "table is currently being updated"
		setSyncedCondition(desired, corev1.ConditionFalse, &msg, nil)
		return desired, requeueWaitWhileUpdating
	}
	if tableHasTerminalStatus(latest) {
		msg := "table is in '" + *latest.ko.Status.TableStatus + "' status"
		setTerminalCondition(desired, corev1.ConditionTrue, &msg, nil)
		setSyncedCondition(desired, corev1.ConditionTrue, nil, nil)
		return desired, nil
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	if delta.DifferentAt("Spec.TimeToLive") {
		if err := rm.syncTTL(ctx, desired, latest); err != nil {
			// Ignore "already disabled errors"
			if awsErr, ok := ackerr.AWSError(err); ok && !(awsErr.Code() == "ValidationException" &&
				strings.HasPrefix(awsErr.Message(), "TimeToLive is already disabled")) {
				return nil, err
			}
		}
	}
	if delta.DifferentAt("Spec.Tags") {
		if err := rm.syncTableTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	}

	if delta.DifferentAt("Spec.BillingMode") ||
		delta.DifferentAt("Spec.SSESpecification") || delta.DifferentAt("Spec.TableClass") {
		if err := rm.syncTable(ctx, desired, delta); err != nil {
			return nil, err
		}
	}

	// We want to update fast fields first
	// Then attributes
	// then GSI
	if delta.DifferentExcept("Spec.Tags", "Spec.TimeToLive") {
		switch {
		case delta.DifferentAt("Spec.StreamSpecification"):
			if err := rm.syncTable(ctx, desired, delta); err != nil {
				return nil, err
			}
		case delta.DifferentAt("Spec.ProvisionedThroughput"):
			if err := rm.syncTableProvisionedThroughput(ctx, desired); err != nil {
				return nil, err
			}

		case delta.DifferentAt("Spec.GlobalSecondaryIndexes"):
			if err := rm.syncTableGlobalSecondaryIndexes(ctx, latest, desired); err != nil {
				return nil, err
			}
		}
	}

	return &resource{ko}, nil
}

// syncTable updates a given table billing mode, stream specification
// or SSE specification.
func (rm *resourceManager) syncTable(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTable")
	defer exit(err)

	input, err := rm.newUpdateTablePayload(ctx, r, delta)
	if err != nil {
		return err
	}

	_, err = rm.sdkapi.UpdateTable(input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateTable", err)
	if err != nil {
		return err
	}
	return nil
}

// newUpdateTablePayload constructs the updateTableInput object.
func (rm *resourceManager) newUpdateTablePayload(
	ctx context.Context,
	r *resource,
	delta *ackcompare.Delta,
) (*svcsdk.UpdateTableInput, error) {
	input := &svcsdk.UpdateTableInput{
		TableName: aws.String(*r.ko.Spec.TableName),
	}

	if delta.DifferentAt("Spec.BillingMode") {
		if r.ko.Spec.BillingMode != nil {
			input.BillingMode = aws.String(*r.ko.Spec.BillingMode)
		} else {
			// set biling mode to the default value `PROVISIONED`
			input.BillingMode = aws.String(svcsdk.BillingModeProvisioned)
		}
	}
	if delta.DifferentAt("Spec.StreamSpecification") {
		if r.ko.Spec.StreamSpecification != nil {
			if r.ko.Spec.StreamSpecification.StreamEnabled != nil {
				input.StreamSpecification = &svcsdk.StreamSpecification{
					StreamEnabled: aws.Bool(*r.ko.Spec.StreamSpecification.StreamEnabled),
				}
				// Only set streamViewType when streamSpefication is enabled and streamViewType is non-nil.
				if *r.ko.Spec.StreamSpecification.StreamEnabled && r.ko.Spec.StreamSpecification.StreamViewType != nil {
					input.StreamSpecification.StreamViewType = aws.String(*r.ko.Spec.StreamSpecification.StreamViewType)
				}
			} else {
				input.StreamSpecification = &svcsdk.StreamSpecification{
					StreamEnabled: aws.Bool(false),
				}
			}
		}
	}
	if delta.DifferentAt("Spec.SSESpecification") {
		if r.ko.Spec.SSESpecification != nil {
			input.SSESpecification = &svcsdk.SSESpecification{}
			if r.ko.Spec.SSESpecification.Enabled != nil {
				input.SSESpecification.Enabled = aws.Bool(*r.ko.Spec.SSESpecification.Enabled)
				if r.ko.Spec.SSESpecification.SSEType != nil {
					r.ko.Spec.SSESpecification.SSEType = aws.String(*r.ko.Spec.SSESpecification.SSEType)
				}
				if r.ko.Spec.SSESpecification.KMSMasterKeyID != nil {
					r.ko.Spec.SSESpecification.KMSMasterKeyID = aws.String(*r.ko.Spec.SSESpecification.KMSMasterKeyID)
				}
			} else {
				input.SSESpecification = &svcsdk.SSESpecification{
					Enabled: aws.Bool(false),
				}
			}
		}
	}

	return input, nil
}

// syncTableProvisionedThroughput updates a given table provisioned throughputs
func (rm *resourceManager) syncTableProvisionedThroughput(
	ctx context.Context,
	r *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTableProvisionedThroughput")
	defer exit(err)

	input := &svcsdk.UpdateTableInput{
		TableName:             aws.String(*r.ko.Spec.TableName),
		ProvisionedThroughput: &svcsdk.ProvisionedThroughput{},
	}
	if r.ko.Spec.ProvisionedThroughput != nil {
		if r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != nil {
			input.ProvisionedThroughput.ReadCapacityUnits = aws.Int64(*r.ko.Spec.ProvisionedThroughput.ReadCapacityUnits)
		} else {
			input.ProvisionedThroughput.ReadCapacityUnits = aws.Int64(0)
		}

		if r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != nil {
			input.ProvisionedThroughput.WriteCapacityUnits = aws.Int64(*r.ko.Spec.ProvisionedThroughput.WriteCapacityUnits)
		} else {
			input.ProvisionedThroughput.WriteCapacityUnits = aws.Int64(0)
		}
	} else {
		input.ProvisionedThroughput.ReadCapacityUnits = aws.Int64(0)
		input.ProvisionedThroughput.WriteCapacityUnits = aws.Int64(0)
	}

	_, err = rm.sdkapi.UpdateTable(input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateTable", err)
	if err != nil {
		return err
	}
	return err
}

// setResourceAdditionalFields will describe the fields that are not return by
// DescribeTable calls
func (rm *resourceManager) setResourceAdditionalFields(
	ctx context.Context,
	ko *v1alpha1.Table,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.setResourceAdditionalFields")
	defer func(err error) { exit(err) }(err)

	if tags, err := rm.getResourceTagsPagesWithContext(ctx, string(*ko.Status.ACKResourceMetadata.ARN)); err != nil {
		return err
	} else {
		ko.Spec.Tags = tags
	}

	if ttlSpec, err := rm.getResourceTTLWithContext(ctx, ko.Spec.TableName); err != nil {
		return err
	} else {
		ko.Spec.TimeToLive = ttlSpec
	}

	return nil
}

func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {

	if len(a.ko.Spec.KeySchema) != len(b.ko.Spec.KeySchema) {
		delta.Add("Spec.KeySchema", a.ko.Spec.KeySchema, b.ko.Spec.KeySchema)
	} else if a.ko.Spec.KeySchema != nil && b.ko.Spec.KeySchema != nil {
		if !equalKeySchemaArrays(a.ko.Spec.KeySchema, b.ko.Spec.KeySchema) {
			delta.Add("Spec.KeySchema", a.ko.Spec.KeySchema, b.ko.Spec.KeySchema)
		}
	}

	if len(a.ko.Spec.AttributeDefinitions) != len(b.ko.Spec.AttributeDefinitions) {
		delta.Add("Spec.AttributeDefinitions", a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions)
	} else if a.ko.Spec.AttributeDefinitions != nil && b.ko.Spec.AttributeDefinitions != nil {
		if !equalAttributeDefinitions(a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions) {
			delta.Add("Spec.AttributeDefinitions", a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions)
		}
	}

	if len(a.ko.Spec.GlobalSecondaryIndexes) != len(b.ko.Spec.GlobalSecondaryIndexes) {
		delta.Add("Spec.GlobalSecondaryIndexes", a.ko.Spec.GlobalSecondaryIndexes, b.ko.Spec.GlobalSecondaryIndexes)
	} else if a.ko.Spec.GlobalSecondaryIndexes != nil && b.ko.Spec.GlobalSecondaryIndexes != nil {
		if !equalGlobalSecondaryIndexesArrays(a.ko.Spec.GlobalSecondaryIndexes, b.ko.Spec.GlobalSecondaryIndexes) {
			delta.Add("Spec.GlobalSecondaryIndexes", a.ko.Spec.GlobalSecondaryIndexes, b.ko.Spec.GlobalSecondaryIndexes)
		}
	}

	// See https://github.com/aws-controllers-k8s/community/issues/1595
	if aws.StringValue(a.ko.Spec.BillingMode) == string(v1alpha1.BillingMode_PAY_PER_REQUEST) {
		a.ko.Spec.ProvisionedThroughput = nil
	}
	if aws.StringValue(b.ko.Spec.BillingMode) == string(v1alpha1.BillingMode_PAY_PER_REQUEST) {
		b.ko.Spec.ProvisionedThroughput = nil
	}

	if len(a.ko.Spec.Tags) != len(b.ko.Spec.Tags) {
		delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
	} else if a.ko.Spec.Tags != nil && b.ko.Spec.Tags != nil {
		if !equalTags(a.ko.Spec.Tags, b.ko.Spec.Tags) {
			delta.Add("Spec.Tags", a.ko.Spec.Tags, b.ko.Spec.Tags)
		}
	}

	if a.ko.Spec.TimeToLive == nil && b.ko.Spec.TimeToLive != nil {
		a.ko.Spec.TimeToLive = &v1alpha1.TimeToLiveSpecification{
			Enabled: &DefaultTTLEnabledValue,
		}
	}
}
