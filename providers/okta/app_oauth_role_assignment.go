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
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type AppOAuthRoleAssignmentGenerator struct {
	OktaService
}

func (g *AppOAuthRoleAssignmentGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	// List all applications
	appList, resp, err := client.ApplicationAPI.ListApplications(ctx).Execute()
	if err != nil {
		return err
	}

	allApps := appList
	for resp.HasNextPage() {
		var nextAppList []okta.ListApplications200ResponseInner
		resp, err = resp.Next(&nextAppList)
		if err != nil {
			return err
		}
		allApps = append(allApps, nextAppList...)
	}

	var resources []terraformutils.Resource

	// For each OAuth app, check if it has role assignments
	for _, app := range allApps {
		appID := extractAppID(app)
		appLabel := extractAppLabel(app)

		if appID == "" {
			continue
		}

		// Only check OAuth apps
		if app.OpenIdConnectApplication == nil {
			continue
		}

		// List client roles for this OAuth app using manual REST API call
		roles, err := g.listClientRoles(appID)
		if err != nil || len(roles) == 0 {
			continue
		}

		// Create a resource for each role assignment
		for _, role := range roles {
			if id, ok := role["id"].(string); ok {
				roleType := ""
				if t, ok := role["type"].(string); ok {
					roleType = t
				}
				resources = append(resources, terraformutils.NewSimpleResource(
					id,
					normalizeResourceName(appID+"_"+appLabel+"_"+roleType),
					"okta_app_oauth_role_assignment",
					"okta",
					[]string{},
				))
			}
		}
	}

	g.Resources = resources
	return nil
}

func (g *AppOAuthRoleAssignmentGenerator) listClientRoles(clientID string) ([]map[string]interface{}, error) {
	ctx, apiSupplement, err := g.APISupplementClient()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("/oauth2/v1/clients/%s/roles", clientID)
	req, err := apiSupplement.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var roles []map[string]interface{}
	_, err = apiSupplement.RequestExecutor.Do(ctx, req, &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
