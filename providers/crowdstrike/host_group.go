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
	"github.com/crowdstrike/gofalcon/falcon/client/host_group"
)

type HostGroupGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike host groups
// Supports filtering by ID: --filter="host_group=id1:id2:id3"
func (g *HostGroupGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific host groups are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "host_group"); shouldFilter {
		// Only import specified host groups
		for _, id := range filteredIDs {
			g.Resources = append(g.Resources, createSimpleResource(
				id,
				"host_group_"+id,
				"crowdstrike_host_group",
				HostGroupAllowEmptyValues,
			))
		}
		return nil
	}

	// Query all host group IDs
	queryParams := host_group.QueryHostGroupsParams{
		Context: ctx,
	}

	queryResp, err := client.HostGroup.QueryHostGroups(&queryParams)
	if err != nil {
		return err
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		// No host groups found
		return nil
	}

	// Get detailed information for all host groups
	getParams := host_group.GetHostGroupsParams{
		Context: ctx,
		Ids:     queryResp.Payload.Resources,
	}

	getResp, err := client.HostGroup.GetHostGroups(&getParams)
	if err != nil {
		return err
	}

	// Create resources for each host group
	for _, hostGroup := range getResp.Payload.Resources {
		resourceID := nilStringValue(hostGroup.ID)
		resourceName := nilStringValue(hostGroup.Name)

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			"host_group_"+resourceName,
			"crowdstrike_host_group",
			HostGroupAllowEmptyValues,
		))
	}

	return nil
}
