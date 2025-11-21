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
	"github.com/crowdstrike/gofalcon/falcon/client/prevention_policies"
)

// PreventionPolicyWindowsGenerator
type PreventionPolicyWindowsGenerator struct {
	CrowdStrikeService
}

func (g *PreventionPolicyWindowsGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Windows", "crowdstrike_prevention_policy_windows")
}

// PreventionPolicyLinuxGenerator
type PreventionPolicyLinuxGenerator struct {
	CrowdStrikeService
}

func (g *PreventionPolicyLinuxGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Linux", "crowdstrike_prevention_policy_linux")
}

// PreventionPolicyMacGenerator
type PreventionPolicyMacGenerator struct {
	CrowdStrikeService
}

func (g *PreventionPolicyMacGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Mac", "crowdstrike_prevention_policy_mac")
}

// Common method to initialize prevention policies for a specific platform
func (g *CrowdStrikeService) initPreventionPoliciesForPlatform(platformName, resourceType string) error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	filter := fmt.Sprintf("platform_name:'%s'", platformName)
	queryParams := prevention_policies.QueryCombinedPreventionPoliciesParams{
		Context: ctx,
		Filter:  &filter,
	}

	resp, err := client.PreventionPolicies.QueryCombinedPreventionPolicies(&queryParams)
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
		resourceName := fmt.Sprintf("prevention_policy_%s_%s", platformName, *policy.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			resourceType,
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// DefaultPreventionPolicyWindowsGenerator
type DefaultPreventionPolicyWindowsGenerator struct {
	CrowdStrikeService
}

func (g *DefaultPreventionPolicyWindowsGenerator) InitResources() error {
	return g.initDefaultPreventionPolicy("Windows", "crowdstrike_default_prevention_policy_windows")
}

// DefaultPreventionPolicyLinuxGenerator
type DefaultPreventionPolicyLinuxGenerator struct {
	CrowdStrikeService
}

func (g *DefaultPreventionPolicyLinuxGenerator) InitResources() error {
	return g.initDefaultPreventionPolicy("Linux", "crowdstrike_default_prevention_policy_linux")
}

// DefaultPreventionPolicyMacGenerator
type DefaultPreventionPolicyMacGenerator struct {
	CrowdStrikeService
}

func (g *DefaultPreventionPolicyMacGenerator) InitResources() error {
	return g.initDefaultPreventionPolicy("Mac", "crowdstrike_default_prevention_policy_mac")
}

// Common method to initialize default prevention policy for a specific platform
func (g *CrowdStrikeService) initDefaultPreventionPolicy(platformName, resourceType string) error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	filter := fmt.Sprintf("platform_name:'%s'+name:'platform_default'", platformName)
	queryParams := prevention_policies.QueryCombinedPreventionPoliciesParams{
		Context: ctx,
		Filter:  &filter,
	}

	resp, err := client.PreventionPolicies.QueryCombinedPreventionPolicies(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := *policy.ID
			resourceName := fmt.Sprintf("default_prevention_policy_%s", platformName)

			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				resourceID,
				resourceName,
				resourceType,
				"crowdstrike",
				[]string{},
			))
			break
		}
	}

	return nil
}

// PreventionPolicyPrecedenceGenerator
type PreventionPolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

func (g *PreventionPolicyPrecedenceGenerator) InitResources() error {
	// Prevention policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("prevention_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_prevention_policy_precedence",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// PreventionPolicyAttachmentGenerator
type PreventionPolicyAttachmentGenerator struct {
	CrowdStrikeService
}

func (g *PreventionPolicyAttachmentGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query all prevention policies across all platforms
	queryParams := prevention_policies.QueryCombinedPreventionPoliciesParams{
		Context: ctx,
	}

	resp, err := client.PreventionPolicies.QueryCombinedPreventionPolicies(&queryParams)
	if err != nil {
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	// Create attachment resources for each policy that has host groups
	for _, policy := range resp.Payload.Resources {
		if policy.Groups != nil && len(policy.Groups) > 0 {
			for _, groupID := range policy.Groups {
				resourceID := fmt.Sprintf("%s_%s", *policy.ID, *groupID)
				resourceName := fmt.Sprintf("attachment_%s_%s", *policy.Name, *groupID)

				g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
					resourceID,
					resourceName,
					"crowdstrike_prevention_policy_attachment",
					"crowdstrike",
					[]string{},
				))
			}
		}
	}

	return nil
}
