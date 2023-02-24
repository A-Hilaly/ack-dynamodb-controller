	if resp.Table.SSEDescription != nil {
		f := &svcapitypes.SSESpecification{}
		if resp.Table.SSEDescription.Status != nil {
			f.Enabled = aws.Bool(*resp.Table.SSEDescription.Status == "ENABLED")
		} else {
			f.Enabled = aws.Bool(false)
		}
		if resp.Table.SSEDescription.SSEType != nil {
			f.SSEType = resp.Table.SSEDescription.SSEType
		}
		if resp.Table.SSEDescription.KMSMasterKeyArn != nil {
			f.KMSMasterKeyID = resp.Table.SSEDescription.KMSMasterKeyArn
		}
		ko.Spec.SSESpecification = f
	} else {
		ko.Spec.SSESpecification = nil
	}
	if resp.Table.TableClassSummary != nil {
		ko.Spec.TableClass = resp.Table.TableClassSummary.TableClass
	} else {
		ko.Spec.TableClass = aws.String("STANDARD")
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