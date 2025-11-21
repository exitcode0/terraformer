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

func (g *HostGroupGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
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

	for _, hostGroup := range getResp.Payload.Resources {
		resourceID := *hostGroup.ID
		resourceName := "host_group_" + *hostGroup.Name

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_host_group",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}
