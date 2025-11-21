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

type ProfileMappingGenerator struct {
	OktaService
}

func (g *ProfileMappingGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// List all profile mappings
	mappings, resp, err := client.ProfileMapping.ListProfileMappings(ctx, &query.Params{})
	if err != nil {
		return err
	}

	allMappings := mappings
	for resp.HasNextPage() {
		var nextMappings []*okta.ProfileMapping
		resp, err = resp.Next(ctx, &nextMappings)
		if err != nil {
			return err
		}
		allMappings = append(allMappings, nextMappings...)
	}

	var resources []terraformutils.Resource
	for _, mapping := range allMappings {
		sourceName := ""
		if mapping.Source != nil && mapping.Source.Name != "" {
			sourceName = mapping.Source.Name
		}
		targetName := ""
		if mapping.Target != nil && mapping.Target.Name != "" {
			targetName = mapping.Target.Name
		}

		resources = append(resources, terraformutils.NewSimpleResource(
			mapping.Id,
			normalizeResourceName(mapping.Id+"_"+sourceName+"_to_"+targetName),
			"okta_profile_mapping",
			"okta",
			[]string{},
		))
	}

	g.Resources = resources
	return nil
}
