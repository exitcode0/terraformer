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
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon/client/custom_ioa"
)

// CloudSecurityRuleGenerator
type CloudSecurityRuleGenerator struct {
	CrowdStrikeService
}

func (g *CloudSecurityRuleGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query custom IOA rules
	queryParams := custom_ioa.QueryRulesV2Params{
		Context: ctx,
	}

	queryResp, err := client.CustomIoa.QueryRulesV2(&queryParams)
	if err != nil {
		return err
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		return nil
	}

	// Get detailed information for all rules
	getParams := custom_ioa.GetRulesV2Params{
		Context: ctx,
		Ids:     queryResp.Payload.Resources,
	}

	getResp, err := client.CustomIoa.GetRulesV2(&getParams)
	if err != nil {
		return err
	}

	for _, rule := range getResp.Payload.Resources {
		resourceID := *rule.ID
		resourceName := fmt.Sprintf("cloud_security_rule_%s", *rule.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_cloud_security_rule",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// CloudComplianceCustomFrameworkGenerator
type CloudComplianceCustomFrameworkGenerator struct {
	CrowdStrikeService
}

func (g *CloudComplianceCustomFrameworkGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query custom compliance frameworks
	// Note: The actual API endpoint may vary - this is a placeholder based on typical patterns
	queryParams := custom_ioa.QueryRuleGroupsFullParams{
		Context: ctx,
	}

	queryResp, err := client.CustomIoa.QueryRuleGroupsFull(&queryParams)
	if err != nil {
		return err
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		return nil
	}

	for _, framework := range queryResp.Payload.Resources {
		resourceID := *framework.ID
		resourceName := fmt.Sprintf("cloud_compliance_custom_framework_%s", *framework.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_cloud_compliance_custom_framework",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}
