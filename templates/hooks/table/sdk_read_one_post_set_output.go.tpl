	// Set billing mode in Spec
	if resp.Table.BillingModeSummary != nil && resp.Table.BillingModeSummary.BillingMode != nil {
		ko.Spec.BillingMode = resp.Table.BillingModeSummary.BillingMode
	} else {
		// Default billingMode value is `PROVISIONED`
		ko.Spec.BillingMode = aws.String(svcsdk.BillingModeProvisioned)
	}
	// Set SSESpecification default
	if resp.Table.SSESpecification != nil && resp.Table.SSESpecification.Enabled != nil {
		ko.Spec.BillingMode = resp.Table.BillingModeSummary.BillingMode
	} else {
		// Default billingMode value is `PROVISIONED`
		ko.Spec.BillingMode = aws.String(svcsdk.BillingModeProvisioned)
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