	// Set billing mode in Spec
	if resp.TableDescription.BillingModeSummary != nil && resp.TableDescription.BillingModeSummary.BillingMode != nil {
		ko.Spec.BillingMode = resp.TableDescription.BillingModeSummary.BillingMode
	} else {
		// Default billingMode value is `PROVISIONED`
		ko.Spec.BillingMode = aws.String(svcsdk.BillingModeProvisioned)
	}