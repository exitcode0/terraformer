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
	"github.com/crowdstrike/gofalcon/falcon/client/content_update_policies"
)

// ContentUpdatePolicyGenerator
type ContentUpdatePolicyGenerator struct {
	CrowdStrikeService
}

func (g *ContentUpdatePolicyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := content_update_policies.QueryCombinedContentUpdatePoliciesParams{
		Context: ctx,
	}

	resp, err := client.ContentUpdatePolicies.QueryCombinedContentUpdatePolicies(&queryParams)
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
		resourceName := fmt.Sprintf("content_update_policy_%s", *policy.Name)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_content_update_policy",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}

// DefaultContentUpdatePolicyGenerator
type DefaultContentUpdatePolicyGenerator struct {
	CrowdStrikeService
}

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
		return err
	}

	if resp.Payload == nil || len(resp.Payload.Resources) == 0 {
		return nil
	}

	for _, policy := range resp.Payload.Resources {
		if policy.Name != nil && *policy.Name == "platform_default" {
			resourceID := *policy.ID
			platformName := policy.PlatformName
			resourceName := fmt.Sprintf("default_content_update_policy_%s", platformName)

			g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
				resourceID,
				resourceName,
				"crowdstrike_default_content_update_policy",
				"crowdstrike",
				[]string{},
			))
		}
	}

	return nil
}

// ContentUpdatePolicyPrecedenceGenerator
type ContentUpdatePolicyPrecedenceGenerator struct {
	CrowdStrikeService
}

func (g *ContentUpdatePolicyPrecedenceGenerator) InitResources() error {
	// Content update policy precedence is a special resource that manages ordering
	// We create one resource per platform
	platforms := []string{"windows", "linux", "mac"}

	for _, platform := range platforms {
		resourceID := fmt.Sprintf("content_update_policy_precedence_%s", platform)
		resourceName := fmt.Sprintf("precedence_%s", platform)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_content_update_policy_precedence",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}
