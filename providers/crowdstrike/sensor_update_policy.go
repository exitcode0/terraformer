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
	"github.com/crowdstrike/gofalcon/falcon/client/sensor_update_policies"
)

// SensorUpdatePolicyGenerator
type SensorUpdatePolicyGenerator struct {
	CrowdStrikeService
}

func (g *SensorUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := sensor_update_policies.QueryCombinedSensorUpdatePoliciesV2Params{
		Context: ctx,
	}

	resp, err := client.SensorUpdatePolicies.QueryCombinedSensorUpdatePoliciesV2(&queryParams)
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
		resourceName := fmt.Sprintf("sensor_update_policy_%s", *policy.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_sensor_update_policy",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// DefaultSensorUpdatePolicyGenerator
type DefaultSensorUpdatePolicyGenerator struct {
	CrowdStrikeService
}

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
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := *policy.ID
			platformName := policy.Settings.SensorVersion
			resourceName := fmt.Sprintf("default_sensor_update_policy_%s", platformName)

			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				resourceID,
				resourceName,
				"crowdstrike_default_sensor_update_policy",
				"crowdstrike",
				[]string{},
			))
		}
	}

	return nil
}

// SensorUpdatePolicyPrecedenceGenerator
type SensorUpdatePolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

func (g *SensorUpdatePolicyPrecedenceGenerator) InitResources() error {
	// Sensor update policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("sensor_update_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_sensor_update_policy_precedence",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// SensorUpdatePolicyAttachmentGenerator
type SensorUpdatePolicyAttachmentGenerator struct {
	CrowdStrikeService
}

func (g *SensorUpdatePolicyAttachmentGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := sensor_update_policies.QueryCombinedSensorUpdatePoliciesV2Params{
		Context: ctx,
	}

	resp, err := client.SensorUpdatePolicies.QueryCombinedSensorUpdatePoliciesV2(&queryParams)
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
					"crowdstrike_sensor_update_policy_host_group_attachment",
					"crowdstrike",
					[]string{},
				))
			}
		}
	}

	return nil
}
