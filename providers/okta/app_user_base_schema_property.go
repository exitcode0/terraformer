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

type AppUserBaseSchemaPropertyGenerator struct {
	OktaService
}

func (g *AppUserBaseSchemaPropertyGenerator) InitResources() error {
	ctx, client, err := g.ClientV5()
	if err != nil {
		return err
	}

	// List all applications
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

	// For each app, get its user schema
	for _, app := range allApps {
		appID := extractAppID(app)
		appLabel := extractAppLabel(app)

		if appID == "" {
			continue
		}

		// Get the user schema for this app - use v2 client for this
		ctxV2, clientV2, err := g.Client()
		if err != nil {
			continue
		}

		userSchema, _, err := clientV2.UserSchema.GetApplicationUserSchema(ctxV2, appID)
		if err != nil {
			// If we can't get the schema for this app, skip it
			continue
		}

		// Process base properties if they exist
		if userSchema.Definitions != nil && userSchema.Definitions.Base != nil && userSchema.Definitions.Base.Properties != nil {
			for propName := range userSchema.Definitions.Base.Properties {
				// The resource ID is app_id/property_name
				resourceID := fmt.Sprintf("%s/%s", appID, propName)

				resources = append(resources, terraformutils.NewSimpleResource(
					resourceID,
					normalizeResourceName(fmt.Sprintf("%s_%s_%s", appID, appLabel, propName)),
					"okta_app_user_base_schema_property",
					"okta",
					[]string{},
				))
			}
		}
	}

	g.Resources = resources
	return nil
}
