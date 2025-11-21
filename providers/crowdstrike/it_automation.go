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
	"github.com/crowdstrike/gofalcon/falcon/client/response_policies"
)

// ItAutomationPolicyGenerator
type ItAutomationPolicyGenerator struct {
	CrowdStrikeService
}

func (g *ItAutomationPolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// IT Automation policies use response policies API
	queryParams := response_policies.QueryCombinedRTResponsePoliciesParams{
		Context: ctx,
	}

	resp, err := client.ResponsePolicies.QueryCombinedRTResponsePolicies(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		// Skip default policies - they have their own resource type
		if policy.Name != nil && *policy.Name == "platform_default" {
			continue
		}

		resourceID := *policy.ID
		resourceName := fmt.Sprintf("it_automation_policy_%s", *policy.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_it_automation_policy",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// ItAutomationDefaultPolicyGenerator
type ItAutomationDefaultPolicyGenerator struct {
	CrowdStrikeService
}

func (g *ItAutomationDefaultPolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query for default IT automation policy
	filter := "name:'platform_default'"
	queryParams := response_policies.QueryCombinedRTResponsePoliciesParams{
		Context: ctx,
		Filter:  &filter,
	}

	resp, err := client.ResponsePolicies.QueryCombinedRTResponsePolicies(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := *policy.ID
			platformName := policy.PlatformName
			resourceName := fmt.Sprintf("it_automation_default_policy_%s", platformName)

			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				resourceID,
				resourceName,
				"crowdstrike_it_automation_default_policy",
				"crowdstrike",
				[]string{},
			))
		}
	}

	return nil
}

// ItAutomationPolicyPrecedenceGenerator
type ItAutomationPolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

func (g *ItAutomationPolicyPrecedenceGenerator) InitResources() error {
	// IT Automation policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("it_automation_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_it_automation_policy_precedence",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// ItAutomationTaskGenerator
type ItAutomationTaskGenerator struct {
	CrowdStrikeService
}

func (g *ItAutomationTaskGenerator) InitResources() error {
	// IT Automation tasks are managed through a separate API
	// This is a placeholder implementation
	// The actual API endpoints would need to be determined from the provider implementation
	return nil
}

// ItAutomationTaskGroupGenerator
type ItAutomationTaskGroupGenerator struct {
	CrowdStrikeService
}

func (g *ItAutomationTaskGroupGenerator) InitResources() error {
	// IT Automation task groups are managed through a separate API
	// This is a placeholder implementation
	// The actual API endpoints would need to be determined from the provider implementation
	return nil
}
