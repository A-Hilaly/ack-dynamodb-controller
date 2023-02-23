package table

import (
	"context"

	"github.com/aws-controllers-k8s/dynamodb-controller/apis/v1alpha1"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/dynamodb"
)

// syncTTL updates a dynamodb table's TimeToLive property.
func (rm *resourceManager) syncTTL(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTTL")
	defer func(err error) { exit(err) }(err)

	ttlSpec := &svcsdk.TimeToLiveSpecification{}
	if desired.ko.Spec.TimeToLive != nil {
		ttlSpec.AttributeName = desired.ko.Spec.TimeToLive.AttributeName
		ttlSpec.Enabled = desired.ko.Spec.TimeToLive.Enabled
	} else {
		// In order to disable the TTL, we can't simply call the
		// `UpdateTimeToLive` method with an empty specification. Instead, we
		// must explicitly set the enabled to false and provide the attribute
		// name of the existing TTL.
		currentAttrName := ""
		if latest.ko.Spec.TimeToLive != nil &&
			latest.ko.Spec.TimeToLive.AttributeName != nil {
			currentAttrName = *latest.ko.Spec.TimeToLive.AttributeName
		}

		ttlSpec.SetAttributeName(currentAttrName)
		ttlSpec.SetEnabled(false)
	}

	_, err = rm.sdkapi.UpdateTimeToLiveWithContext(
		ctx,
		&svcsdk.UpdateTimeToLiveInput{
			TableName:               desired.ko.Spec.TableName,
			TimeToLiveSpecification: ttlSpec,
		},
	)
	rm.metrics.RecordAPICall("UPDATE", "UpdateTimeToLive", err)
	return err
}

// getResourceTTLWithContext queries the table TTL of a given resource.
func (rm *resourceManager) getResourceTTLWithContext(ctx context.Context, tableName *string) (*v1alpha1.TimeToLiveSpecification, error) {
	var err error
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getResourceTTLWithContext")
	defer func(err error) { exit(err) }(err)

	res, err := rm.sdkapi.DescribeTimeToLiveWithContext(
		ctx,
		&svcsdk.DescribeTimeToLiveInput{
			TableName: tableName,
		},
	)
	rm.metrics.RecordAPICall("GET", "DescribeTimeToLive", err)
	if err != nil {
		return nil, err
	}

	// Treat status "ENABLING" and "ENABLED" as `Enabled` == true
	isEnabled := *res.TimeToLiveDescription.TimeToLiveStatus == svcsdk.TimeToLiveStatusEnabled ||
		*res.TimeToLiveDescription.TimeToLiveStatus == svcsdk.TimeToLiveStatusEnabling

	return &v1alpha1.TimeToLiveSpecification{
		AttributeName: res.TimeToLiveDescription.AttributeName,
		Enabled:       &isEnabled,
	}, nil
}
