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

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/terraform-provider-okta/sdk"
)

type AdminRoleCustomGenerator struct {
	OktaService
}

type customRoleResponse struct {
	Roles []struct {
		Id          string `json:"id"`
		Label       string `json:"label"`
		Description string `json:"description"`
	} `json:"roles"`
}

func (g *AdminRoleCustomGenerator) InitResources() error {
	ctx, client, err := g.APISupplementClient()
	if err != nil {
		return err
	}

	// List all custom roles via manual REST API call
	roles, err := listCustomRoles(ctx, client)
	if err != nil {
		return err
	}

	var resources []terraformutils.Resource
	for _, role := range roles {
		resources = append(resources, terraformutils.NewSimpleResource(
			role.Id,
			normalizeResourceName(role.Id+"_"+role.Label),
			"okta_admin_role_custom",
			"okta",
			[]string{},
		))
	}

	g.Resources = resources
	return nil
}

func listCustomRoles(ctx context.Context, m *sdk.APISupplement) ([]struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}, error) {
	url := "/api/v1/iam/roles"
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var response customRoleResponse
	_, err = m.RequestExecutor.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Roles, nil
}
