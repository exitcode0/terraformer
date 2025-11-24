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
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon/client/prevention_policies"
)

// PreventionPolicyWindowsGenerator imports Windows prevention policies
type PreventionPolicyWindowsGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike Windows prevention policies
// Supports filtering by ID: --filter="prevention_policy_windows=id1:id2:id3"
func (g *PreventionPolicyWindowsGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Windows", "crowdstrike_prevention_policy_windows", "prevention_policy_windows")
}

// PreventionPolicyLinuxGenerator imports Linux prevention policies
type PreventionPolicyLinuxGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike Linux prevention policies
// Supports filtering by ID: --filter="prevention_policy_linux=id1:id2:id3"
func (g *PreventionPolicyLinuxGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Linux", "crowdstrike_prevention_policy_linux", "prevention_policy_linux")
}

// PreventionPolicyMacGenerator imports macOS prevention policies
type PreventionPolicyMacGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike macOS prevention policies
// Supports filtering by ID: --filter="prevention_policy_mac=id1:id2:id3"
func (g *PreventionPolicyMacGenerator) InitResources() error {
	return g.initPreventionPoliciesForPlatform("Mac", "crowdstrike_prevention_policy_mac", "prevention_policy_mac")
}

// Common method to initialize prevention policies for a specific platform with filter support
func (g *CrowdStrikeService) initPreventionPoliciesForPlatform(platformName, resourceType, filterKey string) error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific policies are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, filterKey); shouldFilter {
		// Only import specified policies
		for _, id := range filteredIDs {
			g.Resources = append(g.Resources, createSimpleResource(
				id,
				fmt.Sprintf("prevention_policy_%s_%s", platformName, id),
				resourceType,
				PreventionPolicyAllowEmptyValues,
			))
		}
		return nil
	}

	// Query all policies for this platform
	filter := fmt.Sprintf("platform_name:'%s'", platformName)
	queryParams := prevention_policies.QueryCombinedPreventionPoliciesParams{
		Context: ctx,
		Filter:  &filter,
	}

	resp, err := client.PreventionPolicies.QueryCombinedPreventionPolicies(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query %s prevention policies: %w", platformName, err)
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		// Skip default policies - they have their own resource type
		if policy.Name != nil && *policy.Name == "platform_default" {
			continue
		}

		resourceID := nilStringValue(policy.ID)
		resourceName := nilStringValue(policy.Name)

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			fmt.Sprintf("prevention_policy_%s_%s", platformName, resourceName),
			resourceType,
			PreventionPolicyAllowEmptyValues,
		))
	}

	return nil
}

// DefaultPreventionPolicyWindowsGenerator imports default Windows prevention policy
type DefaultPreventionPolicyWindowsGenerator struct {
	CrowdStrikeService
}

// InitResources imports the default Windows prevention policy
func (g *DefaultPreventionPolicyWindowsGenerator) InitResources() error {
	return g.initDefaultPreventionPolicy("Windows", "crowdstrike_default_prevention_policy_windows")
}

// DefaultPreventionPolicyLinuxGenerator imports default Linux prevention policy
type DefaultPreventionPolicyLinuxGenerator struct {
	CrowdStrikeService
}

// InitResources imports the default Linux prevention policy
func (g *DefaultPreventionPolicyLinuxGenerator) InitResources() error {
	return g.initDefaultPreventionPolicy("Linux", "crowdstrike_default_prevention_policy_linux")
}

// DefaultPreventionPolicyMacGenerator imports default macOS prevention policy
type DefaultPreventionPolicyMacGenerator struct {
	CrowdStrikeService
}

// InitResources imports the default macOS prevention policy
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
		return fmt.Errorf("failed to query default %s prevention policy: %w", platformName, err)
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := nilStringValue(policy.ID)
			resourceName := fmt.Sprintf("default_prevention_policy_%s", strings.ToLower(platformName))

			g.Resources = append(g.Resources, createResourceWithName(
				resourceID,
				resourceName,
				resourceType,
				PreventionPolicyAllowEmptyValues,
			))
			break
		}
	}

	return nil
}

// PreventionPolicyPrecedenceGenerator imports prevention policy precedence configurations
type PreventionPolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

// InitResources imports prevention policy precedence (one per platform)
// Supports filtering by platform: --filter="prevention_policy_precedence=windows:linux:mac"
func (g *PreventionPolicyPrecedenceGenerator) InitResources() error {
	// Prevention policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	// Check if specific platforms are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "prevention_policy_precedence"); shouldFilter {
		// Filter platforms by the requested values
		var filteredPlatforms []string
		for _, platform := range platforms {
			for _, filter := range filteredIDs {
				if strings.EqualFold(filter, platform) {
					filteredPlatforms = append(filteredPlatforms, platform)
					break
				}
			}
		}
		platforms = filteredPlatforms
	}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("prevention_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			resourceName,
			"crowdstrike_prevention_policy_precedence",
			PreventionPolicyAllowEmptyValues,
		))
	}

	return nil
}

// PreventionPolicyAttachmentGenerator imports prevention policy-host group attachments
type PreventionPolicyAttachmentGenerator struct {
	CrowdStrikeService
}

// InitResources imports prevention policy attachments to host groups
// Supports filtering by policy ID: --filter="prevention_policy_attachment=policy_id1:policy_id2"
func (g *PreventionPolicyAttachmentGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific policy attachments are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "prevention_policy_attachment"); shouldFilter {
		// Import only attachments for specified policy IDs
		for _, policyID := range filteredIDs {
			// Get specific policy details
			getParams := prevention_policies.GetPreventionPoliciesParams{
				Context: ctx,
				Ids:     []string{policyID},
			}

			policyResp, err := client.PreventionPolicies.GetPreventionPolicies(&getParams)
			if err != nil {
				continue // Skip if policy not found
			}

			if policyResp.Payload != nil && len(policyResp.Payload.Resources) > 0 {
				for _, policy := range policyResp.Payload.Resources {
					g.addPolicyAttachments(policy)
				}
			}
		}
		return nil
	}

	// Query all prevention policies across all platforms
	queryParams := prevention_policies.QueryCombinedPreventionPoliciesParams{
		Context: ctx,
	}

	resp, err := client.PreventionPolicies.QueryCombinedPreventionPolicies(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query prevention policies for attachments: %w", err)
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	// Create attachment resources for each policy that has host groups
	for _, policy := range resp.Payload.Resources {
		g.addPolicyAttachments(policy)
	}

	return nil
}

// addPolicyAttachments creates resources for policy-host group attachments
func (g *PreventionPolicyAttachmentGenerator) addPolicyAttachments(policy interface{}) {
	// Type assertion to get policy details
	policyMap, ok := policy.(map[string]interface{})
	if !ok {
		return
	}

	policyID, _ := policyMap["id"].(string)
	policyName, _ := policyMap["name"].(string)
	groups, _ := policyMap["groups"].([]interface{})

	if len(groups) > 0 {
		for _, group := range groups {
			groupID, _ := group.(string)
			if groupID != "" {
				resourceID := fmt.Sprintf("%s_%s", policyID, groupID)
				resourceName := sanitizeResourceName(fmt.Sprintf("attachment_%s_%s", policyName, groupID))

				g.Resources = append(g.Resources, createResourceWithName(
					resourceID,
					resourceName,
					"crowdstrike_prevention_policy_attachment",
					PreventionPolicyAllowEmptyValues,
				))
			}
		}
	}
}
