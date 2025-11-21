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

type GroupOwnerGenerator struct {
	OktaService
}

func (g *GroupOwnerGenerator) InitResources() error {
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

	// For each group, list owners using v5 SDK
	ctxV5, clientV5, err := g.ClientV5()
	if err != nil {
		return err
	}

	for _, group := range groupList {
		groupID := group.Id
		groupName := group.Profile.Name

		// List group owners
		owners, _, err := clientV5.GroupOwnerAPI.ListGroupOwners(ctxV5, groupID).Execute()
		if err != nil {
			// If we can't list owners for this group, skip it
			continue
		}

		// Create a resource for each owner
		for _, owner := range owners {
			ownerID := owner.GetId()
			resources = append(resources, terraformutils.NewSimpleResource(
				ownerID,
				normalizeResourceName(fmt.Sprintf("%s_%s_%s", groupID, groupName, ownerID)),
				"okta_group_owner",
				"okta",
				[]string{},
			))
		}
	}

	g.Resources = resources
	return nil
}
