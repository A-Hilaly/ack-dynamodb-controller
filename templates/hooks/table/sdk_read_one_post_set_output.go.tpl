	// Set billing mode in Spec
	if resp.Table.BillingModeSummary != nil && resp.Table.BillingModeSummary.BillingMode != nil {
		ko.Spec.BillingMode = resp.Table.BillingModeSummary.BillingMode
	} else {
		// Default billingMode value is `PROVISIONED`
		ko.Spec.BillingMode = aws.String(svcsdk.BillingModeProvisioned)
	}
	// Set SSESpecification
	if resp.Table.SSEDescription != nil {
		sseSpecification := &svcapitypes.SSESpecification{}
		if resp.Table.SSEDescription.Status != nil {
			switch *resp.Table.SSEDescription.Status {
			case string(svcapitypes.SSEStatus_ENABLED),
				string(svcapitypes.SSEStatus_ENABLING),
				string(svcapitypes.SSEStatus_UPDATING):

				sseSpecification.Enabled = aws.Bool(true)
				if resp.Table.SSEDescription.SSEType != nil {
					sseSpecification.SSEType = resp.Table.SSEDescription.SSEType
				}
			case string(svcapitypes.SSEStatus_DISABLED),
				string(svcapitypes.SSEStatus_DISABLING):
				sseSpecification.Enabled = aws.Bool(false)
			}
		}
		ko.Spec.SSESpecification = sseSpecification
	} else {
		ko.Spec.SSESpecification = &svcapitypes.SSESpecification{
			Enabled: aws.Bool(false),
		}
	}
	if isTableCreating(&resource{ko}) {
		return &resource{ko}, requeueWaitWhileCreating
	}
	if isTableUpdating(&resource{ko}) {
		return &resource{ko}, requeueWaitWhileUpdating
	}
	if err := rm.setResourceAdditionalFields(ctx, ko); err != nil {
		return nil, err
	}