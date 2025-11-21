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
)

type GroupSchemaPropertyGenerator struct {
	OktaService
}

func (g *GroupSchemaPropertyGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Get group schema
	groupSchema, _, err := client.GroupSchema.GetGroupSchema(ctx)
	if err != nil {
		return err
	}

	var resources []terraformutils.Resource

	// Process custom properties
	if groupSchema.Definitions != nil && groupSchema.Definitions.Custom != nil && groupSchema.Definitions.Custom.Properties != nil {
		for key := range groupSchema.Definitions.Custom.Properties {
			resources = append(resources, terraformutils.NewSimpleResource(
				key,
				normalizeResourceName("group_custom_"+key),
				"okta_group_schema_property",
				"okta",
				[]string{},
			))
		}
	}

	g.Resources = resources
	return nil
}
