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

type RequestConditionGenerator struct {
	OktaService
}

func (g *RequestConditionGenerator) InitResources() error {
	ctx, client, err := g.APISupplementClient()
	if err != nil {
		return err
	}

	// List all request conditions via manual REST API call
	// This is a governance feature
	conditions, err := listRequestConditions(ctx, client)
	if err != nil {
		return err
	}

	var resources []terraformutils.Resource
	for _, condition := range conditions {
		// The resource ID is resource_id/condition_id
		resourceID := fmt.Sprintf("%s/%s", condition.ResourceID, condition.ID)

		resources = append(resources, terraformutils.NewSimpleResource(
			resourceID,
			normalizeResourceName(condition.ResourceID+"_"+condition.Name),
			"okta_request_condition",
			"okta",
			[]string{},
		))
	}

	g.Resources = resources
	return nil
}

type requestCondition struct {
	ID         string
	ResourceID string
	Name       string
}

func listRequestConditions(ctx context.Context, m *sdk.APISupplement) ([]requestCondition, error) {
	// Request conditions are part of Okta Governance
	// First get all resources, then get conditions for each
	url := "/api/v1/governance/resources"
	req, err := m.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var governanceResources []map[string]interface{}
	_, err = m.RequestExecutor.Do(ctx, req, &governanceResources)
	if err != nil {
		// Governance might not be enabled, return empty list
		return []requestCondition{}, nil
	}

	result := make([]requestCondition, 0)

	// For each governance resource, list its conditions
	for _, resource := range governanceResources {
		resourceID, ok := resource["id"].(string)
		if !ok {
			continue
		}

		conditionsURL := fmt.Sprintf("/api/v1/governance/resources/%s/conditions", resourceID)
		condReq, err := m.RequestExecutor.NewRequest("GET", conditionsURL, nil)
		if err != nil {
			continue
		}

		var conditions []map[string]interface{}
		_, err = m.RequestExecutor.Do(ctx, condReq, &conditions)
		if err != nil {
			continue
		}

		for _, cond := range conditions {
			if id, ok := cond["id"].(string); ok {
				condition := requestCondition{
					ID:         id,
					ResourceID: resourceID,
				}
				if name, ok := cond["name"].(string); ok {
					condition.Name = name
				}
				result = append(result, condition)
			}
		}
	}

	return result, nil
}
