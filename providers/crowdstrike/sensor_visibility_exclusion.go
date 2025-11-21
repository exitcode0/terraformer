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

// SensorVisibilityExclusionGenerator
type SensorVisibilityExclusionGenerator struct {
	CrowdStrikeService
}

func (g *SensorVisibilityExclusionGenerator) InitResources() error {
	ctx, client, err := g.Client()
	if err != nil {
		return err
	}

	queryParams := sensor_visibility_exclusions.QuerySensorVisibilityExclusionsV1Params{
		Context: ctx,
	}

	queryResp, err := client.SensorVisibilityExclusions.QuerySensorVisibilityExclusionsV1(&queryParams)
	if err != nil {
		return err
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
		return err
	}

	for _, exclusion := range getResp.Payload.Resources {
		resourceID := *exclusion.ID
		resourceName := fmt.Sprintf("sensor_visibility_exclusion_%s", *exclusion.Value)

		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			resourceID,
			resourceName,
			"crowdstrike_sensor_visibility_exclusion",
			"crowdstrike",
			[]string{},
		))
	}

	return nil
}
