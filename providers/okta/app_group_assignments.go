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

type AppGroupAssignmentsGenerator struct {
	OktaService
}

func (g *AppGroupAssignmentsGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	// First, list all applications
	appList, resp, err := client.ApplicationAPI.ListApplications(ctx).Execute()
	if err != nil {
		return fmt.Errorf("error listing applications: %w", err)
	}

	allApps := appList
	for resp.HasNextPage() {
		var nextAppList []okta.ListApplications200ResponseInner
		resp, err = resp.Next(&nextAppList)
		if err != nil {
			return fmt.Errorf("error fetching next page of applications: %w", err)
		}
		allApps = append(allApps, nextAppList...)
	}

	var resources []terraformutils.Resource

	// For each app, check if it has group assignments
	for _, app := range allApps {
		appID := extractAppID(app)
		appLabel := extractAppLabel(app)

		if appID == "" {
			continue
		}

		// List group assignments for this app - using v2 client as v5 doesn't have this method
		ctxV2, clientV2, err := g.Client()
		if err != nil {
			continue
		}

		groupAssignments, _, err := clientV2.Application.ListApplicationGroupAssignments(ctxV2, appID, nil)
		if err != nil {
			// If we can't list assignments for this app, skip it
			continue
		}

		// Only create a resource if there are group assignments
		if len(groupAssignments) > 0 {
			resources = append(resources, terraformutils.NewSimpleResource(
				appID,
				normalizeResourceName(appID+"_"+appLabel),
				"okta_app_group_assignments",
				"okta",
				[]string{},
			))
		}
	}

	g.Resources = resources
	return nil
}

func extractAppID(app okta.ListApplications200ResponseInner) string {
	if app.OpenIdConnectApplication != nil && app.OpenIdConnectApplication.Id != nil {
		return *app.OpenIdConnectApplication.Id
	}
	if app.SamlApplication != nil && app.SamlApplication.Id != nil {
		return *app.SamlApplication.Id
	}
	if app.Saml11Application != nil && app.Saml11Application.Id != nil {
		return *app.Saml11Application.Id
	}
	if app.BookmarkApplication != nil && app.BookmarkApplication.Id != nil {
		return *app.BookmarkApplication.Id
	}
	if app.AutoLoginApplication != nil && app.AutoLoginApplication.Id != nil {
		return *app.AutoLoginApplication.Id
	}
	if app.BasicAuthApplication != nil && app.BasicAuthApplication.Id != nil {
		return *app.BasicAuthApplication.Id
	}
	if app.BrowserPluginApplication != nil && app.BrowserPluginApplication.Id != nil {
		return *app.BrowserPluginApplication.Id
	}
	if app.SecurePasswordStoreApplication != nil && app.SecurePasswordStoreApplication.Id != nil {
		return *app.SecurePasswordStoreApplication.Id
	}
	return ""
}

func extractAppLabel(app okta.ListApplications200ResponseInner) string {
	if app.OpenIdConnectApplication != nil {
		return app.OpenIdConnectApplication.Label
	}
	if app.SamlApplication != nil {
		return app.SamlApplication.Label
	}
	if app.Saml11Application != nil {
		return app.Saml11Application.Label
	}
	if app.BookmarkApplication != nil {
		return app.BookmarkApplication.Label
	}
	if app.AutoLoginApplication != nil {
		return app.AutoLoginApplication.Label
	}
	if app.BasicAuthApplication != nil {
		return app.BasicAuthApplication.Label
	}
	if app.BrowserPluginApplication != nil {
		return app.BrowserPluginApplication.Label
	}
	if app.SecurePasswordStoreApplication != nil {
		return app.SecurePasswordStoreApplication.Label
	}
	return ""
}
