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
	"github.com/crowdstrike/gofalcon/falcon/client/filevantage"
)

// FilevantagePolicyGenerator
type FilevantagePolicyGenerator struct {
	CrowdStrikeService
}

func (g *FilevantagePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := filevantage.QueryPoliciesParams{
		Context: ctx,
	}

	queryResp, err := client.Filevantage.QueryPolicies(&queryParams)
	if err != nil {
		return err
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		return nil
	}

	// Get detailed information for all policies
	getParams := filevantage.GetPoliciesParams{
		Context: ctx,
		Ids:     queryResp.Payload.Resources,
	}

	getResp, err := client.Filevantage.GetPolicies(&getParams)
	if err != nil {
		return err
	}

	for _, policy := range getResp.Payload.Resources {
		resourceID := *policy.ID
		resourceName := fmt.Sprintf("filevantage_policy_%s", *policy.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_filevantage_policy",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// FilevantageRuleGroupGenerator
type FilevantageRuleGroupGenerator struct {
	CrowdStrikeService
}

func (g *FilevantageRuleGroupGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := filevantage.QueryRuleGroupsParams{
		Context: ctx,
	}

	queryResp, err := client.Filevantage.QueryRuleGroups(&queryParams)
	if err != nil {
		return err
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		return nil
	}

	// Get detailed information for all rule groups
	getParams := filevantage.GetRuleGroupsParams{
		Context: ctx,
		Ids:     queryResp.Payload.Resources,
	}

	getResp, err := client.Filevantage.GetRuleGroups(&getParams)
	if err != nil {
		return err
	}

	for _, ruleGroup := range getResp.Payload.Resources {
		resourceID := *ruleGroup.ID
		resourceName := fmt.Sprintf("filevantage_rule_group_%s", *ruleGroup.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_filevantage_rule_group",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// FilevantagePolicyPrecedenceGenerator
type FilevantagePolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

func (g *FilevantagePolicyPrecedenceGenerator) InitResources() error {
	// FileVantage policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("filevantage_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_filevantage_policy_precedence",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}
