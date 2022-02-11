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
	"reflect"
	"testing"

	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
)

func Test_computeGlobalSecondaryIndexDelta(t *testing.T) {
	type args struct {
		a []*v1alpha1.GlobalSecondaryIndex
		b []*v1alpha1.GlobalSecondaryIndex
	}
	tests := []struct {
		name        string
		args        args
		wantAdded   []*v1alpha1.GlobalSecondaryIndex
		wantUpdated []*v1alpha1.GlobalSecondaryIndex
		wantRemoved []string
	}{
		{
			name: "nil global secondary index arrays",
			args: args{
				a: nil,
				b: nil,
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: nil,
		},
		{
			name: "nil and empty global secondary index arrays",
			args: args{
				a: nil,
				b: []*v1alpha1.GlobalSecondaryIndex{},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: nil,
		},
		{
			name: "empty global secondary index arrays",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{},
				b: []*v1alpha1.GlobalSecondaryIndex{},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: nil,
		},
		{
			name: "One added global secondary array",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
				},
			},
			wantAdded: []*v1alpha1.GlobalSecondaryIndex{
				{IndexName: aws.String("new_index_1")},
			},
			wantRemoved: nil,
			wantUpdated: nil,
		},
		{
			name: "Multiple added global secondary array",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
					{IndexName: aws.String("new_index_2")},
					{IndexName: aws.String("new_index_3")},
				},
			},
			wantAdded: []*v1alpha1.GlobalSecondaryIndex{
				{IndexName: aws.String("new_index_2")},
				{IndexName: aws.String("new_index_3")},
			},
			wantRemoved: nil,
			wantUpdated: nil,
		},
		{
			name: "One removed global secondary array",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{},
			},
			wantAdded: nil,
			wantRemoved: []string{
				"new_index_1",
			},
			wantUpdated: nil,
		},
		{
			name: "Multiple removed global secondary array",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
					{IndexName: aws.String("new_index_2")},
					{IndexName: aws.String("new_index_3")},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{IndexName: aws.String("new_index_1")},
				},
			},
			wantAdded: nil,
			wantRemoved: []string{
				"new_index_2",
				"new_index_3",
			},
			wantUpdated: nil,
		},
		{
			name: "Updated global secondary array - Projection.ProjectType",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						Projection: &v1alpha1.Projection{
							ProjectionType: aws.String("type1"),
						},
					},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						Projection: &v1alpha1.Projection{
							ProjectionType: aws.String("type2"),
						},
					},
				},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_1"),
					Projection: &v1alpha1.Projection{
						ProjectionType: aws.String("type2"),
					},
				},
			},
		},
		{
			name: "Updated global secondary array - Projection.NonKeyAttributes",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						Projection: &v1alpha1.Projection{
							NonKeyAttributes: []*string{aws.String("k1")},
						},
					},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						Projection: &v1alpha1.Projection{
							NonKeyAttributes: []*string{aws.String("k2")},
						},
					},
				},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_1"),
					Projection: &v1alpha1.Projection{
						NonKeyAttributes: []*string{aws.String("k2")},
					},
				},
			},
		},
		{
			name: "Updated global secondary array - ProvisionedThroughput",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						ProvisionedThroughput: &v1alpha1.ProvisionedThroughput{
							ReadCapacityUnits:  aws.Int64(10),
							WriteCapacityUnits: aws.Int64(20),
						},
					},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						ProvisionedThroughput: &v1alpha1.ProvisionedThroughput{
							ReadCapacityUnits:  aws.Int64(10),
							WriteCapacityUnits: aws.Int64(30),
						},
					},
				},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_1"),
					ProvisionedThroughput: &v1alpha1.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(10),
						WriteCapacityUnits: aws.Int64(30),
					},
				},
			},
		},
		{
			name: "Updated global secondary array - KeySchema",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						KeySchema: []*v1alpha1.KeySchemaElement{
							{
								AttributeName: aws.String("attr1"),
								KeyType:       aws.String("integer"),
							},
						},
					},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						KeySchema: []*v1alpha1.KeySchemaElement{
							{
								AttributeName: aws.String("attr1"),
								KeyType:       aws.String("integer"),
							},
							{
								AttributeName: aws.String("attr2"),
								KeyType:       aws.String("integer"),
							},
						},
					},
				},
			},
			wantAdded:   nil,
			wantRemoved: nil,
			wantUpdated: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_1"),
					KeySchema: []*v1alpha1.KeySchemaElement{
						{
							AttributeName: aws.String("attr1"),
							KeyType:       aws.String("integer"),
						},
						{
							AttributeName: aws.String("attr2"),
							KeyType:       aws.String("integer"),
						},
					},
				},
			},
		},
		{
			name: "Added, removed and updated global secondary indexes",
			args: args{
				a: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						KeySchema: []*v1alpha1.KeySchemaElement{
							{
								AttributeName: aws.String("attr1"),
								KeyType:       aws.String("integer"),
							},
						},
					},
					{
						IndexName: aws.String("new_index_2"),
					},
				},
				b: []*v1alpha1.GlobalSecondaryIndex{
					{
						IndexName: aws.String("new_index_1"),
						KeySchema: []*v1alpha1.KeySchemaElement{
							{
								AttributeName: aws.String("attr1"),
								KeyType:       aws.String("integer"),
							},
							{
								AttributeName: aws.String("attr2"),
								KeyType:       aws.String("integer"),
							},
						},
					},
					{
						IndexName: aws.String("new_index_3"),
						ProvisionedThroughput: &v1alpha1.ProvisionedThroughput{
							ReadCapacityUnits:  aws.Int64(10),
							WriteCapacityUnits: aws.Int64(20),
						},
					},
				},
			},
			wantAdded: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_3"),
					ProvisionedThroughput: &v1alpha1.ProvisionedThroughput{
						ReadCapacityUnits:  aws.Int64(10),
						WriteCapacityUnits: aws.Int64(20),
					},
				},
			},
			wantRemoved: []string{"new_index_2"},
			wantUpdated: []*v1alpha1.GlobalSecondaryIndex{
				{
					IndexName: aws.String("new_index_1"),
					KeySchema: []*v1alpha1.KeySchemaElement{
						{
							AttributeName: aws.String("attr1"),
							KeyType:       aws.String("integer"),
						},
						{
							AttributeName: aws.String("attr2"),
							KeyType:       aws.String("integer"),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAdded, gotUpdated, gotRemoved := computeGlobalSecondaryIndexDelta(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(gotAdded, tt.wantAdded) {
				t.Errorf("computeGlobalSecondaryIndexDelta() gotAdded = %v, want %v", gotAdded, tt.wantAdded)
			}
			if !reflect.DeepEqual(gotUpdated, tt.wantUpdated) {
				t.Errorf("computeGlobalSecondaryIndexDelta() gotUpdated = %v, want %v", gotUpdated, tt.wantUpdated)
			}
			if !reflect.DeepEqual(gotRemoved, tt.wantRemoved) {
				t.Errorf("computeGlobalSecondaryIndexDelta() gotRemoved = %v, want %v", gotRemoved, tt.wantRemoved)
			}
		})
	}
}
