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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"

	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
)

func computeGlobalSecondaryIndexDelta(
	a []*v1alpha1.GlobalSecondaryIndex,
	b []*v1alpha1.GlobalSecondaryIndex,
) (added, updated []*v1alpha1.GlobalSecondaryIndex, removed []string) {
	var visitedIndexes []string
	for _, aElement := range a {
		visitedIndexes = append(visitedIndexes, *aElement.IndexName)
		found := false
		for _, bElement := range b {
			if *aElement.IndexName == *bElement.IndexName {
				found = true
				if !equalGlobalSecondaryIndexes(aElement, bElement) {
					updated = append(updated, aElement)
				}
			}
			if !found {
				removed = append(removed, *aElement.IndexName)
			}
		}

		for _, bElement := range b {
			if !ackutil.InStrings(*bElement.IndexName, visitedIndexes) {
				added = append(added, bElement)
			}
		}
	}
	return added, updated, removed
}

func equalKeySchemas(
	a []*v1alpha1.KeySchemaElement,
	b []*v1alpha1.KeySchemaElement,
) bool {
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if equalStrings(aElement.AttributeName, bElement.AttributeName) {
				found = true
				if !equalStrings(aElement.KeyType, bElement.KeyType) {
					return false
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func equalGlobalSecondaryIndexes(
	a *v1alpha1.GlobalSecondaryIndex,
	b *v1alpha1.GlobalSecondaryIndex,
) bool {
	if ackcompare.HasNilDifference(a.ProvisionedThroughput, b.ProvisionedThroughput) {
		return false
	}
	if a.ProvisionedThroughput != nil && b.ProvisionedThroughput != nil {
		if !equalInt64(a.ProvisionedThroughput.ReadCapacityUnits, b.ProvisionedThroughput.ReadCapacityUnits) {
			return false
		}
		if equalInt64(a.ProvisionedThroughput.WriteCapacityUnits, b.ProvisionedThroughput.WriteCapacityUnits) {
			return false
		}
	}
	if ackcompare.HasNilDifference(a.Projection, b.Projection) {
		return false
	}
	if a.Projection != nil && b.Projection != nil {
		if !equalStrings(a.Projection.ProjectionType, b.Projection.ProjectionType) {
			return false
		}
		if !ackcompare.SliceStringPEqual(a.Projection.NonKeyAttributes, b.Projection.NonKeyAttributes) {
			return false
		}
	}
	return true
}

func customPreCompare(
	delta *ackcompare.Delta,
	a *resource,
	b *resource,
) {
	if ackcompare.HasNilDifference(a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions) ||
		len(a.ko.Spec.AttributeDefinitions) != len(b.ko.Spec.AttributeDefinitions) {
		delta.Add("Spec.AttributeDefinitions", a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions)
	} else if a.ko.Spec.AttributeDefinitions != nil && b.ko.Spec.AttributeDefinitions != nil {
		if !equalAttributeDefinitions(a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions) {
			delta.Add("Spec.AttributeDefinitions", a.ko.Spec.AttributeDefinitions, b.ko.Spec.AttributeDefinitions)
		}
	}
}

func equalAttributeDefinitions(
	a []*v1alpha1.AttributeDefinition,
	b []*v1alpha1.AttributeDefinition,
) bool {
	for _, aElement := range a {
		found := false
		for _, bElement := range b {
			if equalStrings(aElement.AttributeName, bElement.AttributeName) {
				found = true
				if !equalStrings(aElement.AttributeType, bElement.AttributeType) {
					return false
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func emptyString(s *string) bool {
	if s == nil {
		return true
	}
	return *s == ""
}

func equalStrings(a, b *string) bool {
	if a == nil {
		return b == nil || *b == ""
	}
	return (*a == "" && b == nil) || *a == *b
}

func equalInt64(a, b *int64) bool {
	if a == nil {
		return b == nil || *b == 0
	}
	return (*a == 0 && b == nil) || *a == *b
}
