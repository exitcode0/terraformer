// Copyright 2018 The Terraformer Authors.
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

package crowdstrike

import (
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon/client/sensor_visibility_exclusions"
)

// SensorVisibilityExclusionGenerator imports sensor visibility exclusions
type SensorVisibilityExclusionGenerator struct {
	CrowdStrikeService
}

// InitResources imports CrowdStrike sensor visibility exclusions
// Supports filtering by ID: --filter="sensor_visibility_exclusion=id1:id2:id3"
func (g *SensorVisibilityExclusionGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	// Check if specific exclusions are requested via filter
	if filteredIDs, shouldFilter := getFilteredIDs(&g.Service, "sensor_visibility_exclusion"); shouldFilter {
		// Only import specified exclusions
		for _, id := range filteredIDs {
			g.Resources = append(g.Resources, createSimpleResource(
				id,
				fmt.Sprintf("exclusion_%s", id),
				"crowdstrike_sensor_visibility_exclusion",
				[]string{},
			))
		}
		return nil
	}

	queryParams := sensor_visibility_exclusions.QuerySensorVisibilityExclusionsV1Params{
		Context: ctx,
	}

	queryResp, err := client.SensorVisibilityExclusions.QuerySensorVisibilityExclusionsV1(&queryParams)
	if err != nil {
		return fmt.Errorf("failed to query sensor visibility exclusions: %w", err)
	}

	if queryResp.Payload == nil || len(queryResp.Payload.Resources) == 0 {
		return nil
	}

	// Get detailed information for all exclusions
	getParams := sensor_visibility_exclusions.GetSensorVisibilityExclusionsV1Params{
		Context: ctx,
		Ids:     queryResp.Payload.Resources,
	}

	getResp, err := client.SensorVisibilityExclusions.GetSensorVisibilityExclusionsV1(&getParams)
	if err != nil {
		return fmt.Errorf("failed to get sensor visibility exclusion details: %w", err)
	}

	for _, exclusion := range getResp.Payload.Resources {
		resourceID := nilStringValue(exclusion.ID)
		resourceValue := nilStringValue(exclusion.Value)
		resourceName := sanitizeResourceName(fmt.Sprintf("exclusion_%s", resourceValue))

		g.Resources = append(g.Resources, createResourceWithName(
			resourceID,
			resourceName,
			"crowdstrike_sensor_visibility_exclusion",
			[]string{},
		))
	}

	return nil
}
