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

func computeLocalSecondaryIndexDelta(
	a []*v1alpha1.LocalSecondaryIndex,
	b []*v1alpha1.LocalSecondaryIndex,
) (added, updated []*v1alpha1.LocalSecondaryIndex, removed []string) {
	var visitedIndexes []string
loopA:
	for _, aElement := range a {
		visitedIndexes = append(visitedIndexes, *aElement.IndexName)
		for _, bElement := range b {
			if *aElement.IndexName == *bElement.IndexName {
				if !equalLocalSecondaryIndexes(aElement, bElement) {
					updated = append(updated, bElement)
				}
				continue loopA
			}
		}
		removed = append(removed, *aElement.IndexName)

	}
	for _, bElement := range b {
		if !ackutil.InStrings(*bElement.IndexName, visitedIndexes) {
			added = append(added, bElement)
		}
	}
	return added, updated, removed
}

// equalTags returns true if two LocalSecondaryIndex arrays are equal regardless
// of the order of their elements.
func equalLocalSecondaryIndexesArrays(
	a []*v1alpha1.LocalSecondaryIndex,
	b []*v1alpha1.LocalSecondaryIndex,
) bool {
	added, updated, removed := computeLocalSecondaryIndexDelta(a, b)
	return len(added) == 0 && len(updated) == 0 && len(removed) == 0
}

// equalLocalSecondaryIndexes returns whether two LocalSecondaryIndex objects are
// equal or not.
func equalLocalSecondaryIndexes(
	a *v1alpha1.LocalSecondaryIndex,
	b *v1alpha1.LocalSecondaryIndex,
) bool {
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
	if len(a.KeySchema) != len(b.KeySchema) {
		return false
	} else if len(a.KeySchema) > 0 {
		if !equalKeySchemaArrays(a.KeySchema, b.KeySchema) {
			return false
		}
	}
	return true
}
