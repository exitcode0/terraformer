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
	"github.com/crowdstrike/gofalcon/falcon/client/content_update_policies"
)

// ContentUpdatePolicyGenerator imports content update policies
type ContentUpdatePolicyGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike content update policies
// Supports filtering by ID: --filter="content_update_policy=id1:id2:id3"
func (g *ContentUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific policies are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "content_update_policy"); shouldFilter {
		// Only import specified policies
		for _, id := range filteredIDs {
			g.Resources = append(g.Resources, createSimpleResource(
				id,
				fmt.Sprintf("content_update_policy_%s", id),
				"crowdstrike_content_update_policy",
				ContentUpdatePolicyAllowEmptyValues,
			))
		}
		return nil
	}

	queryParams := content_update_policies.QueryCombinedContentUpdatePoliciesParams{
		Context: ctx,
	}

	resp, err := client.ContentUpdatePolicies.QueryCombinedContentUpdatePolicies(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query content update policies: %w", err)
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
			fmt.Sprintf("content_update_policy_%s", resourceName),
			"crowdstrike_content_update_policy",
			ContentUpdatePolicyAllowEmptyValues,
		))
	}

	return nil
}

// DefaultContentUpdatePolicyGenerator imports default content update policy
type DefaultContentUpdatePolicyGenerator struct {
	CrowdStrikeService
}

// InitResources imports the default content update policy
func (g *DefaultContentUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query for default policies across all platforms
	filter := "name:'platform_default'"
	queryParams := content_update_policies.QueryCombinedContentUpdatePoliciesParams{
		Context: ctx,
		Filter:  &filter,
	}

	resp, err := client.ContentUpdatePolicies.QueryCombinedContentUpdatePolicies(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query default content update policy: %w", err)
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := nilStringValue(policy.ID)
			platformName := nilStringValue(policy.PlatformName)
			resourceName := fmt.Sprintf("default_content_update_policy_%s", strings.ToLower(platformName))

			g.Resources = append(g.Resources, createResourceWithName(
				resourceID,
				resourceName,
				"crowdstrike_default_content_update_policy",
				ContentUpdatePolicyAllowEmptyValues,
			))
		}
	}

	return nil
}

// ContentUpdatePolicyPrecedenceGenerator imports content update policy precedence configurations
type ContentUpdatePolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

// InitResources imports content update policy precedence (one per platform)
// Supports filtering by platform: --filter="content_update_policy_precedence=windows:linux:mac"
func (g *ContentUpdatePolicyPrecedenceGenerator) InitResources() error {
	// Content update policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	// Check if specific platforms are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "content_update_policy_precedence"); shouldFilter {
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
		resourceID := fmt.Sprintf("content_update_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			resourceName,
			"crowdstrike_content_update_policy_precedence",
			ContentUpdatePolicyAllowEmptyValues,
		))
	}

	return nil
}
