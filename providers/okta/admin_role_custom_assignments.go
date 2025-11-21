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
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/terraform-provider-okta/sdk"
)

type AdminRoleCustomAssignmentsGenerator struct {
	OktaService
}

type bindingsResponse struct {
	Bindings []struct {
		Role string `json:"role"`
	} `json:"bindings"`
}

func (g *AdminRoleCustomAssignmentsGenerator) InitResources() error {
	ctx, apiSupplement, err := g.APISupplementClient()
	if err != nil {
		return err
	}

	// First, list all resource sets
	resourceSets, err := listResourceSets(ctx, apiSupplement)
	if err != nil {
		return err
	}

	var resources []terraformutils.Resource

	// For each resource set, list bindings (role assignments)
	for _, rs := range resourceSets {
		bindings, err := listResourceSetBindings(ctx, apiSupplement, rs.Id)
		if err != nil {
			// If we can't list bindings for this resource set, skip it
			continue
		}

		// Create a resource for each binding
		for _, binding := range bindings {
			// The resource ID is resource_set_id/custom_role_id
			resourceID := fmt.Sprintf("%s/%s", rs.Id, binding.Role)

			resources = append(resources, terraformutils.NewSimpleResource(
				resourceID,
				normalizeResourceName(fmt.Sprintf("%s_%s_%s", rs.Id, rs.Label, binding.Role)),
				"okta_admin_role_custom_assignments",
				"okta",
				[]string{},
			))
		}
	}

	g.Resources = resources
	return nil
}

func listResourceSetBindings(ctx context.Context, m *sdk.APISupplement, resourceSetID string) ([]struct {
	Role string `json:"role"`
}, error) {
	url := fmt.Sprintf("/api/v1/iam/resource-sets/%s/bindings", resourceSetID)
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var response bindingsResponse
	_, err = m.RequestExecutor.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Bindings, nil
}
