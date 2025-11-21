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

// Generic generator for all device assurance policies
type PolicyDeviceAssuranceGenerator struct {
	OktaService
	Platform string // WINDOWS, MACOS, IOS, ANDROID, CHROMEOS
}

func (g *PolicyDeviceAssuranceGenerator) InitResources() error {
	ctx, apiSupplement, err := g.APISupplementClient()
	if err != nil {
		return err
	}

	policies, err := g.listDeviceAssurancePolicies(ctx, apiSupplement)
	if err != nil {
		return err
	}

	var resources []terraformutils.Resource
	for _, policy := range policies {
		// Filter by platform if specified
		if g.Platform != "" && policy.Platform != g.Platform {
			continue
		}

		resourceType := "okta_policy_device_assurance_" + normalizeResourceName(policy.Platform)

		resources = append(resources, terraformutils.NewSimpleResource(
			policy.ID,
			normalizeResourceName(policy.ID+"_"+policy.Name),
			resourceType,
			"okta",
			[]string{},
		))
	}

	g.Resources = resources
	return nil
}

type deviceAssurancePolicy struct {
	ID       string
	Name     string
	Platform string
}

func (g *PolicyDeviceAssuranceGenerator) listDeviceAssurancePolicies(ctx context.Context, apiSupplement *sdk.APISupplement) ([]deviceAssurancePolicy, error) {
	// Manual REST API call as device assurance is not in all SDK versions
	url := "/api/v1/device-assurances"
	req, err := apiSupplement.RequestExecutor.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var policies []map[string]interface{}
	_, err = apiSupplement.RequestExecutor.Do(ctx, req, &policies)
	if err != nil {
		return nil, err
	}

	result := make([]deviceAssurancePolicy, 0)
	for _, p := range policies {
		if id, ok := p["id"].(string); ok {
			policy := deviceAssurancePolicy{
				ID: id,
			}
			if name, ok := p["name"].(string); ok {
				policy.Name = name
			}
			if platform, ok := p["platform"].(string); ok {
				policy.Platform = platform
			}
			result = append(result, policy)
		}
	}

	return result, nil
}

// Specific generators for each platform
type PolicyDeviceAssuranceWindowsGenerator struct {
	PolicyDeviceAssuranceGenerator
}

func (g *PolicyDeviceAssuranceWindowsGenerator) InitResources() error {
	g.Platform = "WINDOWS"
	return g.PolicyDeviceAssuranceGenerator.InitResources()
}

type PolicyDeviceAssuranceMacOSGenerator struct {
	PolicyDeviceAssuranceGenerator
}

func (g *PolicyDeviceAssuranceMacOSGenerator) InitResources() error {
	g.Platform = "MACOS"
	return g.PolicyDeviceAssuranceGenerator.InitResources()
}

type PolicyDeviceAssuranceIOSGenerator struct {
	PolicyDeviceAssuranceGenerator
}

func (g *PolicyDeviceAssuranceIOSGenerator) InitResources() error {
	g.Platform = "IOS"
	return g.PolicyDeviceAssuranceGenerator.InitResources()
}

type PolicyDeviceAssuranceAndroidGenerator struct {
	PolicyDeviceAssuranceGenerator
}

func (g *PolicyDeviceAssuranceAndroidGenerator) InitResources() error {
	g.Platform = "ANDROID"
	return g.PolicyDeviceAssuranceGenerator.InitResources()
}

type PolicyDeviceAssuranceChromeOSGenerator struct {
	PolicyDeviceAssuranceGenerator
}

func (g *PolicyDeviceAssuranceChromeOSGenerator) InitResources() error {
	g.Platform = "CHROMEOS"
	return g.PolicyDeviceAssuranceGenerator.InitResources()
}
