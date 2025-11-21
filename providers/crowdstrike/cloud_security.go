// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crowdstrike

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon/client/cspm_registration"
)

// CloudAWSAccountGenerator
type CloudAWSAccountGenerator struct {
	CrowdStrikeService
}

func (g *CloudAWSAccountGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := cspm_registration.GetCSPMAwsAccountParams{
		Context: ctx,
	}

	resp, err := client.CspmRegistration.GetCSPMAwsAccount(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, account := range resp.Payload.Resources {
		resourceID := account.AccountID
		resourceName := "aws_account_" + account.AccountID

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_cloud_aws_account",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// CloudAzureTenantGenerator
type CloudAzureTenantGenerator struct {
	CrowdStrikeService
}

func (g *CloudAzureTenantGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := cspm_registration.GetCSPMAzureTenantIDsParams{
		Context: ctx,
	}

	resp, err := client.CspmRegistration.GetCSPMAzureTenantIDs(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	// Get detailed information for each tenant
	for _, tenantID := range resp.Payload.Resources {
		// Query for the tenant details
		getParams := cspm_registration.GetCSPMAzureAccountParams{
			Context: ctx,
			Ids:     []string{tenantID},
		}

		getResp, err := client.CspmRegistration.GetCSPMAzureAccount(&getParams)
		if err != nil {
			continue
		}

		if getResp.Payload == nil || len(getResp.Payload.Resources) == 0 {
			continue
		}

		for _, tenant := range getResp.Payload.Resources {
			resourceID := tenant.TenantID
			resourceName := "azure_tenant_" + tenant.TenantID

			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				resourceID,
				resourceName,
				"crowdstrike_cloud_azure_tenant",
				"crowdstrike",
				[]string{},
			))
		}
	}

	return nil
}
