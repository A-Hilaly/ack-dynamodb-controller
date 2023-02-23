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
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	acktags "github.com/aws-controllers-k8s/runtime/pkg/tags"
)

// Hack to avoid import errors during build...
var (
	_ = &bytes.Buffer{}
	_ = &reflect.Method{}
	_ = &acktags.Tags{}
)

// newResourceDelta returns a new `ackcompare.Delta` used to compare two
// resources
func newResourceDelta(
	a *resource,
	b *resource,
) *ackcompare.Delta {
	delta := ackcompare.NewDelta()
	if (a == nil && b != nil) ||
		(a != nil && b == nil) {
		delta.Add("", a, b)
		return delta
	}
	customPreCompare(delta, a, b)

	if ackcompare.HasNilDifference(a.ko.Spec.BillingMode, b.ko.Spec.BillingMode) {
		delta.Add("Spec.BillingMode", a.ko.Spec.BillingMode, b.ko.Spec.BillingMode)
	} else if a.ko.Spec.BillingMode != nil && b.ko.Spec.BillingMode != nil {
		if *a.ko.Spec.BillingMode != *b.ko.Spec.BillingMode {
			delta.Add("Spec.BillingMode", a.ko.Spec.BillingMode, b.ko.Spec.BillingMode)
		}
	}
	if !reflect.DeepEqual(a.ko.Spec.LocalSecondaryIndexes, b.ko.Spec.LocalSecondaryIndexes) {
		delta.Add("Spec.LocalSecondaryIndexes", a.ko.Spec.LocalSecondaryIndexes, b.ko.Spec.LocalSecondaryIndexes)
	}
	if ackcompare.HasNilDifference(a.ko.Spec.ProvisionedThroughput, b.ko.Spec.ProvisionedThroughput) {
		delta.Add("Spec.ProvisionedThroughput", a.ko.Spec.ProvisionedThroughput, b.ko.Spec.ProvisionedThroughput)
	} else if a.ko.Spec.ProvisionedThroughput != nil && b.ko.Spec.ProvisionedThroughput != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.ProvisionedThroughput.ReadCapacityUnits, b.ko.Spec.ProvisionedThroughput.ReadCapacityUnits) {
			delta.Add("Spec.ProvisionedThroughput.ReadCapacityUnits", a.ko.Spec.ProvisionedThroughput.ReadCapacityUnits, b.ko.Spec.ProvisionedThroughput.ReadCapacityUnits)
		} else if a.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != nil && b.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != nil {
			if *a.ko.Spec.ProvisionedThroughput.ReadCapacityUnits != *b.ko.Spec.ProvisionedThroughput.ReadCapacityUnits {
				delta.Add("Spec.ProvisionedThroughput.ReadCapacityUnits", a.ko.Spec.ProvisionedThroughput.ReadCapacityUnits, b.ko.Spec.ProvisionedThroughput.ReadCapacityUnits)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.ProvisionedThroughput.WriteCapacityUnits, b.ko.Spec.ProvisionedThroughput.WriteCapacityUnits) {
			delta.Add("Spec.ProvisionedThroughput.WriteCapacityUnits", a.ko.Spec.ProvisionedThroughput.WriteCapacityUnits, b.ko.Spec.ProvisionedThroughput.WriteCapacityUnits)
		} else if a.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != nil && b.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != nil {
			if *a.ko.Spec.ProvisionedThroughput.WriteCapacityUnits != *b.ko.Spec.ProvisionedThroughput.WriteCapacityUnits {
				delta.Add("Spec.ProvisionedThroughput.WriteCapacityUnits", a.ko.Spec.ProvisionedThroughput.WriteCapacityUnits, b.ko.Spec.ProvisionedThroughput.WriteCapacityUnits)
			}
		}
	}
	fmt.Println("STARTING SSE DELTA++")
	x, _ := json.Marshal(a.ko.Spec.SSESpecification)
	fmt.Println("SSE VALUE A", string(x))
	x, _ = json.Marshal(b.ko.Spec.SSESpecification)
	fmt.Println("SSE VALUE B", string(x))
	if ackcompare.HasNilDifference(a.ko.Spec.SSESpecification, b.ko.Spec.SSESpecification) {
		delta.Add("Spec.SSESpecification", a.ko.Spec.SSESpecification, b.ko.Spec.SSESpecification)
	} else if a.ko.Spec.SSESpecification != nil && b.ko.Spec.SSESpecification != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.SSESpecification.Enabled, b.ko.Spec.SSESpecification.Enabled) {
			delta.Add("Spec.SSESpecification.Enabled", a.ko.Spec.SSESpecification.Enabled, b.ko.Spec.SSESpecification.Enabled)
		} else if a.ko.Spec.SSESpecification.Enabled != nil && b.ko.Spec.SSESpecification.Enabled != nil {
			if *a.ko.Spec.SSESpecification.Enabled != *b.ko.Spec.SSESpecification.Enabled {
				delta.Add("Spec.SSESpecification.Enabled", a.ko.Spec.SSESpecification.Enabled, b.ko.Spec.SSESpecification.Enabled)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.SSESpecification.KMSMasterKeyID, b.ko.Spec.SSESpecification.KMSMasterKeyID) {
			delta.Add("Spec.SSESpecification.KMSMasterKeyID", a.ko.Spec.SSESpecification.KMSMasterKeyID, b.ko.Spec.SSESpecification.KMSMasterKeyID)
		} else if a.ko.Spec.SSESpecification.KMSMasterKeyID != nil && b.ko.Spec.SSESpecification.KMSMasterKeyID != nil {
			if *a.ko.Spec.SSESpecification.KMSMasterKeyID != *b.ko.Spec.SSESpecification.KMSMasterKeyID {
				delta.Add("Spec.SSESpecification.KMSMasterKeyID", a.ko.Spec.SSESpecification.KMSMasterKeyID, b.ko.Spec.SSESpecification.KMSMasterKeyID)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.SSESpecification.SSEType, b.ko.Spec.SSESpecification.SSEType) {
			delta.Add("Spec.SSESpecification.SSEType", a.ko.Spec.SSESpecification.SSEType, b.ko.Spec.SSESpecification.SSEType)
		} else if a.ko.Spec.SSESpecification.SSEType != nil && b.ko.Spec.SSESpecification.SSEType != nil {
			if *a.ko.Spec.SSESpecification.SSEType != *b.ko.Spec.SSESpecification.SSEType {
				delta.Add("Spec.SSESpecification.SSEType", a.ko.Spec.SSESpecification.SSEType, b.ko.Spec.SSESpecification.SSEType)
			}
		}
	}
	fmt.Println("ENDING SSE DELTA", delta.Differences)

	if ackcompare.HasNilDifference(a.ko.Spec.StreamSpecification, b.ko.Spec.StreamSpecification) {
		delta.Add("Spec.StreamSpecification", a.ko.Spec.StreamSpecification, b.ko.Spec.StreamSpecification)
	} else if a.ko.Spec.StreamSpecification != nil && b.ko.Spec.StreamSpecification != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.StreamSpecification.StreamEnabled, b.ko.Spec.StreamSpecification.StreamEnabled) {
			delta.Add("Spec.StreamSpecification.StreamEnabled", a.ko.Spec.StreamSpecification.StreamEnabled, b.ko.Spec.StreamSpecification.StreamEnabled)
		} else if a.ko.Spec.StreamSpecification.StreamEnabled != nil && b.ko.Spec.StreamSpecification.StreamEnabled != nil {
			if *a.ko.Spec.StreamSpecification.StreamEnabled != *b.ko.Spec.StreamSpecification.StreamEnabled {
				delta.Add("Spec.StreamSpecification.StreamEnabled", a.ko.Spec.StreamSpecification.StreamEnabled, b.ko.Spec.StreamSpecification.StreamEnabled)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.StreamSpecification.StreamViewType, b.ko.Spec.StreamSpecification.StreamViewType) {
			delta.Add("Spec.StreamSpecification.StreamViewType", a.ko.Spec.StreamSpecification.StreamViewType, b.ko.Spec.StreamSpecification.StreamViewType)
		} else if a.ko.Spec.StreamSpecification.StreamViewType != nil && b.ko.Spec.StreamSpecification.StreamViewType != nil {
			if *a.ko.Spec.StreamSpecification.StreamViewType != *b.ko.Spec.StreamSpecification.StreamViewType {
				delta.Add("Spec.StreamSpecification.StreamViewType", a.ko.Spec.StreamSpecification.StreamViewType, b.ko.Spec.StreamSpecification.StreamViewType)
			}
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.TableClass, b.ko.Spec.TableClass) {
		delta.Add("Spec.TableClass", a.ko.Spec.TableClass, b.ko.Spec.TableClass)
	} else if a.ko.Spec.TableClass != nil && b.ko.Spec.TableClass != nil {
		if *a.ko.Spec.TableClass != *b.ko.Spec.TableClass {
			delta.Add("Spec.TableClass", a.ko.Spec.TableClass, b.ko.Spec.TableClass)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.TableName, b.ko.Spec.TableName) {
		delta.Add("Spec.TableName", a.ko.Spec.TableName, b.ko.Spec.TableName)
	} else if a.ko.Spec.TableName != nil && b.ko.Spec.TableName != nil {
		if *a.ko.Spec.TableName != *b.ko.Spec.TableName {
			delta.Add("Spec.TableName", a.ko.Spec.TableName, b.ko.Spec.TableName)
		}
	}
	if ackcompare.HasNilDifference(a.ko.Spec.TimeToLive, b.ko.Spec.TimeToLive) {
		delta.Add("Spec.TimeToLive", a.ko.Spec.TimeToLive, b.ko.Spec.TimeToLive)
	} else if a.ko.Spec.TimeToLive != nil && b.ko.Spec.TimeToLive != nil {
		if ackcompare.HasNilDifference(a.ko.Spec.TimeToLive.AttributeName, b.ko.Spec.TimeToLive.AttributeName) {
			delta.Add("Spec.TimeToLive.AttributeName", a.ko.Spec.TimeToLive.AttributeName, b.ko.Spec.TimeToLive.AttributeName)
		} else if a.ko.Spec.TimeToLive.AttributeName != nil && b.ko.Spec.TimeToLive.AttributeName != nil {
			if *a.ko.Spec.TimeToLive.AttributeName != *b.ko.Spec.TimeToLive.AttributeName {
				delta.Add("Spec.TimeToLive.AttributeName", a.ko.Spec.TimeToLive.AttributeName, b.ko.Spec.TimeToLive.AttributeName)
			}
		}
		if ackcompare.HasNilDifference(a.ko.Spec.TimeToLive.Enabled, b.ko.Spec.TimeToLive.Enabled) {
			delta.Add("Spec.TimeToLive.Enabled", a.ko.Spec.TimeToLive.Enabled, b.ko.Spec.TimeToLive.Enabled)
		} else if a.ko.Spec.TimeToLive.Enabled != nil && b.ko.Spec.TimeToLive.Enabled != nil {
			if *a.ko.Spec.TimeToLive.Enabled != *b.ko.Spec.TimeToLive.Enabled {
				delta.Add("Spec.TimeToLive.Enabled", a.ko.Spec.TimeToLive.Enabled, b.ko.Spec.TimeToLive.Enabled)
			}
		}
	}

	return delta
}
