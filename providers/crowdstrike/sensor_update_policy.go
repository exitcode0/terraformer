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
	"github.com/crowdstrike/gofalcon/falcon/client/sensor_update_policies"
)

// SensorUpdatePolicyGenerator imports sensor update policies
type SensorUpdatePolicyGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike sensor update policies
// Supports filtering by ID: --filter="sensor_update_policy=id1:id2:id3"
func (g *SensorUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific policies are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "sensor_update_policy"); shouldFilter {
		// Only import specified policies
		for _, id := range filteredIDs {
			g.Resources = append(g.Resources, createSimpleResource(
				id,
				fmt.Sprintf("sensor_update_policy_%s", id),
				"crowdstrike_sensor_update_policy",
				SensorUpdatePolicyAllowEmptyValues,
			))
		}
		return nil
	}

	// Query all sensor update policies
	queryParams := sensor_update_policies.QueryCombinedSensorUpdatePoliciesV2Params{
		Context: ctx,
	}

	resp, err := client.SensorUpdatePolicies.QueryCombinedSensorUpdatePoliciesV2(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query sensor update policies: %w", err)
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
			fmt.Sprintf("sensor_update_policy_%s", resourceName),
			"crowdstrike_sensor_update_policy",
			SensorUpdatePolicyAllowEmptyValues,
		))
	}

	return nil
}

// DefaultSensorUpdatePolicyGenerator imports default sensor update policy
type DefaultSensorUpdatePolicyGenerator struct {
	CrowdStrikeService
}

// InitResources imports the default sensor update policy
func (g *DefaultSensorUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Query for default policies across all platforms
	name := "platform_default"
	queryParams := sensor_update_policies.QueryCombinedSensorUpdatePoliciesV2Params{
		Context: ctx,
		Filter:  &name,
	}

	resp, err := client.SensorUpdatePolicies.QueryCombinedSensorUpdatePoliciesV2(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query default sensor update policy: %w", err)
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := nilStringValue(policy.ID)
			platformName := nilStringValue(policy.Settings.SensorVersion)
			resourceName := fmt.Sprintf("default_sensor_update_policy_%s", platformName)

			g.Resources = append(g.Resources, createResourceWithName(
				resourceID,
				resourceName,
				"crowdstrike_default_sensor_update_policy",
				SensorUpdatePolicyAllowEmptyValues,
			))
		}
	}

	return nil
}

// SensorUpdatePolicyPrecedenceGenerator imports sensor update policy precedence configurations
type SensorUpdatePolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

// InitResources imports sensor update policy precedence (one per platform)
// Supports filtering by platform: --filter="sensor_update_policy_precedence=windows:linux:mac"
func (g *SensorUpdatePolicyPrecedenceGenerator) InitResources() error {
	// Sensor update policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	// Check if specific platforms are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "sensor_update_policy_precedence"); shouldFilter {
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
		resourceID := fmt.Sprintf("sensor_update_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			resourceName,
			"crowdstrike_sensor_update_policy_precedence",
			SensorUpdatePolicyAllowEmptyValues,
		))
	}

	return nil
}

// SensorUpdatePolicyAttachmentGenerator imports sensor update policy-host group attachments
type SensorUpdatePolicyAttachmentGenerator struct {
	CrowdStrikeService
}

// InitResources imports sensor update policy attachments to host groups
// Supports filtering by policy ID: --filter="sensor_update_policy_host_group_attachment=policy_id1:policy_id2"
func (g *SensorUpdatePolicyAttachmentGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific policy attachments are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "sensor_update_policy_host_group_attachment"); shouldFilter {
		// Import only attachments for specified policy IDs
		for _, policyID := range filteredIDs {
			// Get specific policy details
			getParams := sensor_update_policies.GetSensorUpdatePoliciesV2Params{
				Context: ctx,
				Ids:     []string{policyID},
			}

			policyResp, err := client.SensorUpdatePolicies.GetSensorUpdatePoliciesV2(&getParams)
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

	queryParams := sensor_update_policies.QueryCombinedSensorUpdatePoliciesV2Params{
		Context: ctx,
	}

	resp, err := client.SensorUpdatePolicies.QueryCombinedSensorUpdatePoliciesV2(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query sensor update policies for attachments: %w", err)
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
func (g *SensorUpdatePolicyAttachmentGenerator) addPolicyAttachments(policy interface{}) {
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
					"crowdstrike_sensor_update_policy_host_group_attachment",
					SensorUpdatePolicyAllowEmptyValues,
				))
			}
		}
	}
}
