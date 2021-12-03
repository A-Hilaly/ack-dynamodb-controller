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
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"
)

func debug(i interface{}) {
	bytes, _ := json.Marshal(i)
	fmt.Println(string(bytes))
}

func debugAll(is ...interface{}) {
	for _, i := range is {
		debug(i)
	}
}

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
				if !reflect.DeepEqual(aElement, bElement) {
					updated = append(updated, bElement)
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

/* found = true
	if ackcompare.HasNilDifference(aElements.ProvisionedThroughput, bElements.ProvisionedThroughput) {
		updated = append(updated, bElements)
	} else if aElements.ProvisionedThroughput != nil && bElements.ProvisionedThroughput != nil {
		if ackcompare.HasNilDifference(aElements.ProvisionedThroughput.ReadCapacityUnits, bElements.ProvisionedThroughput.ReadCapacityUnits) {
			updated = append(updated, bElements)
		} else if aElements.ProvisionedThroughput.ReadCapacityUnits != nil && bElements.ProvisionedThroughput.ReadCapacityUnits != nil {
			if *aElements.ProvisionedThroughput.ReadCapacityUnits != *bElements.ProvisionedThroughput.ReadCapacityUnits {
				updated = append(updated, bElements)
			}
		}
		if ackcompare.HasNilDifference(aElements.ProvisionedThroughput.WriteCapacityUnits, bElements.ProvisionedThroughput.WriteCapacityUnits) {
			updated = append(updated, bElements)
		} else if aElements.ProvisionedThroughput.WriteCapacityUnits != nil && bElements.ProvisionedThroughput.WriteCapacityUnits != nil {
			if *aElements.ProvisionedThroughput.WriteCapacityUnits != *bElements.ProvisionedThroughput.WriteCapacityUnits {
				updated = append(updated, bElements)
			}
		}
	}
	if ackcompare.HasNilDifference(
		aElements.Projection,
		bElements.Projection,
	) {
		updated = append(updated, bElements)
	} else if aElements.Projection. != nil && bElements.ProvisionedThroughput.WriteCapacityUnits != nil {
		if *aElements.ProvisionedThroughput.WriteCapacityUnits != *bElements.ProvisionedThroughput.WriteCapacityUnits {
			updated = append(updated, bElements)
		}
	}
} */
//updated = append(updated, bElements)
