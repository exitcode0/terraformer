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
	"github.com/okta/okta-sdk-golang/v5/okta"
)

type AppOAuthAPIScopeGenerator struct {
	OktaService
}

func (g *AppOAuthAPIScopeGenerator) InitResources() error {
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

	// For each OAuth app, check if it has scope consent grants
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

		// Use v2 client to list scope consent grants
		ctxV2, clientV2, err := g.Client()
		if err != nil {
			continue
		}

		scopes, _, err := clientV2.Application.ListScopeConsentGrants(ctxV2, appID, nil)
		if err != nil || len(scopes) == 0 {
			// If we can't list scopes or there are none, skip it
			continue
		}

		// Create a resource for this app's API scopes
		resources = append(resources, terraformutils.NewSimpleResource(
			appID,
			normalizeResourceName(appID+"_"+appLabel+"_api_scopes"),
			"okta_app_oauth_api_scope",
			"okta",
			[]string{},
		))
	}

	g.Resources = resources
	return nil
}
