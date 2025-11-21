// Copyright 2019 The Terraformer Authors.
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

package okta

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type GroupRoleGenerator struct {
	OktaService
}

func (g *GroupRoleGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// List all groups
	filter := query.NewQueryParams(query.WithFilter("type eq \"OKTA_GROUP\""))
	groupList, resp, err := client.Group.ListGroups(ctx, filter)
	if err != nil {
		return err
	}

	for resp.HasNextPage() {
		var nextGroupSet []*okta.Group
		resp, _ = resp.Next(ctx, &nextGroupSet)
		groupList = append(groupList, nextGroupSet...)
	}

	var resources []terraformutils.Resource

	// For each group, list assigned roles
	for _, group := range groupList {
		groupID := group.Id
		groupName := group.Profile.Name

		// List group assigned roles
		roles, roleResp, err := client.Group.ListGroupAssignedRoles(ctx, groupID, &query.Params{})
		if err != nil {
			// If we can't list roles for this group, skip it
			continue
		}

		// Create a resource for each role
		for _, role := range roles {
			// The resource ID is group_id/role_id
			resourceID := fmt.Sprintf("%s/%s", groupID, role.Id)

			resources = append(resources, terraformutils.NewSimpleResource(
				resourceID,
				normalizeResourceName(fmt.Sprintf("%s_%s_%s", groupID, groupName, role.Type)),
				"okta_group_role",
				"okta",
				[]string{},
			))
		}

		// Handle pagination
		for roleResp.HasNextPage() {
			var nextRoles []*okta.Role
			roleResp, _ = roleResp.Next(ctx, &nextRoles)
			for _, role := range nextRoles {
				resourceID := fmt.Sprintf("%s/%s", groupID, role.Id)
				resources = append(resources, terraformutils.NewSimpleResource(
					resourceID,
					normalizeResourceName(fmt.Sprintf("%s_%s_%s", groupID, groupName, role.Type)),
					"okta_group_role",
					"okta",
					[]string{},
				))
			}
		}
	}

	g.Resources = resources
	return nil
}
