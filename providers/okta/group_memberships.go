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
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type GroupMembershipsGenerator struct {
	OktaService
}

func (g *GroupMembershipsGenerator) InitResources() error {
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

	// For each group, check if it has members
	for _, group := range groupList {
		groupID := group.Id
		groupName := group.Profile.Name

		// List group members
		users, userResp, err := client.Group.ListGroupUsers(ctx, groupID, &query.Params{Limit: 200})
		if err != nil {
			// If we can't list users for this group, skip it
			continue
		}

		// Only create a resource if there are members
		if len(users) > 0 {
			resources = append(resources, terraformutils.NewSimpleResource(
				groupID,
				normalizeResourceName(groupID+"_"+groupName),
				"okta_group_memberships",
				"okta",
				[]string{},
			))
		}

		// Handle pagination
		for userResp.HasNextPage() {
			var nextUsers []*okta.User
			userResp, _ = userResp.Next(ctx, &nextUsers)
		}
	}

	g.Resources = resources
	return nil
}
