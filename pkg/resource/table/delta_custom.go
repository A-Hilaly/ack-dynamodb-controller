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
	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
)

func computeGlobalSecondaryIndexDelta(
	gsia []*v1alpha1.GlobalSecondaryIndex,
	gsib []*v1alpha1.GlobalSecondaryIndex,
) (added, removed, updated []*v1alpha1.GlobalSecondaryIndex) {
	var visitedIndexes []string
	for _, gsiaElem := range gsia {
		visitedIndexes = append(visitedIndexes, *gsiaElem.IndexName)
		found := false
		for _, gsibElem := range gsib {
			if *gsiaElem.IndexName == *gsibElem.IndexName {
				if ackcompare.HasNilDifference(gsiaElem.ProvisionedThroughput, gsibElem.ProvisionedThroughput) {
					updated = append(updated, gsibElem)
				} else if gsiaElem.ProvisionedThroughput != nil && gsibElem.ProvisionedThroughput != nil {
					if ackcompare.HasNilDifference(gsiaElem.ProvisionedThroughput.ReadCapacityUnits, gsibElem.ProvisionedThroughput.ReadCapacityUnits) {
						updated = append(updated, gsibElem)
					} else if gsiaElem.ProvisionedThroughput.ReadCapacityUnits != nil && gsibElem.ProvisionedThroughput.ReadCapacityUnits != nil {
						if *gsiaElem.ProvisionedThroughput.ReadCapacityUnits != *gsibElem.ProvisionedThroughput.ReadCapacityUnits {
							updated = append(updated, gsibElem)
						}
					}
					if ackcompare.HasNilDifference(gsiaElem.ProvisionedThroughput.WriteCapacityUnits, gsibElem.ProvisionedThroughput.WriteCapacityUnits) {
						updated = append(updated, gsibElem)
					} else if gsiaElem.ProvisionedThroughput.WriteCapacityUnits != nil && gsibElem.ProvisionedThroughput.WriteCapacityUnits != nil {
						if *gsiaElem.ProvisionedThroughput.WriteCapacityUnits != *gsibElem.ProvisionedThroughput.WriteCapacityUnits {
							updated = append(updated, gsibElem)
						}
					}
				}
				if ackcompare.HasNilDifference(
					gsibElem.ProvisionedThroughput,
					gsibElem.ProvisionedThroughput,
				) {

				} else {

				}
			}
			updated = append(updated, gsibElem)
		}
		if !found {
			removed = append(added, gsiaElem)
		}
	}

	for _, gsibElem := range gsib {
		_ = gsibElem
	}
	return
}
