package table

import (
	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

// equalKeySchemaArrays return whether two KeySchemaElement arrays are equal or not.
func equalKeySchemaArrays(
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

func equalAttributeDefinitions(a, b []*v1alpha1.AttributeDefinition) bool {
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

func newSDKAttributesDefinition(ads []*v1alpha1.AttributeDefinition) []*svcsdk.AttributeDefinition {
	attributeDefintions := []*svcsdk.AttributeDefinition{}
	for _, ad := range ads {
		attributeDefintion := &svcsdk.AttributeDefinition{}
		if ad != nil {
			if ad.AttributeName != nil {
				attributeDefintion.AttributeName = aws.String(*ad.AttributeName)
			} else {
				attributeDefintion.AttributeName = aws.String("")
			}
			if ad.AttributeType != nil {
				attributeDefintion.AttributeType = aws.String(*ad.AttributeType)
			} else {
				attributeDefintion.AttributeType = aws.String("")
			}
		} else {
			attributeDefintion.AttributeType = aws.String("")
			attributeDefintion.AttributeName = aws.String("")
		}
		attributeDefintions = append(attributeDefintions, attributeDefintion)
	}
	return attributeDefintions
}
