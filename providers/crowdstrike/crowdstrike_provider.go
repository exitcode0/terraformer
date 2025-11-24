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
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client"
	"github.com/zclconf/go-cty/cty"
)

type CrowdStrikeProvider struct { //nolint
	terraformutils.Provider
	clientID     string
	clientSecret string
	cloud        string
	memberCID    string
	client       *client.CrowdStrikeAPISpecification
}

// Init initializes the provider with credentials and creates the Falcon client
func (p *CrowdStrikeProvider) Init(args []string) error {
	// Parse arguments: [clientID, clientSecret, cloud, memberCID]
	// Priority: CLI args > Environment variables

	if len(args) >= 1 && args[0] != "" {
		p.clientID = args[0]
	} else if clientID := os.Getenv("FALCON_CLIENT_ID"); clientID != "" {
		p.clientID = clientID
	} else {
		return errors.New("FALCON_CLIENT_ID is required. Set via --client-id flag or FALCON_CLIENT_ID environment variable")
	}

	if len(args) >= 2 && args[1] != "" {
		p.clientSecret = args[1]
	} else if clientSecret := os.Getenv("FALCON_CLIENT_SECRET"); clientSecret != "" {
		p.clientSecret = clientSecret
	} else {
		return errors.New("FALCON_CLIENT_SECRET is required. Set via --client-secret flag or FALCON_CLIENT_SECRET environment variable")
	}

	if len(args) >= 3 && args[2] != "" {
		p.cloud = args[2]
	} else if cloud := os.Getenv("FALCON_CLOUD"); cloud != "" {
		p.cloud = cloud
	} else {
		p.cloud = "autodiscover" // Default value
	}

	if len(args) >= 4 && args[3] != "" {
		p.memberCID = args[3]
	} else if memberCID := os.Getenv("FALCON_MEMBER_CID"); memberCID != "" {
		p.memberCID = memberCID
	}
	// memberCID is optional for MSSP scenarios

	// Initialize the CrowdStrike Falcon client once
	apiConfig := falcon.ApiConfig{
		ClientId:     p.clientID,
		ClientSecret: p.clientSecret,
		Cloud:        falcon.Cloud(p.cloud),
		Context:      context.Background(),
	}

	if p.memberCID != "" {
		apiConfig.MemberCID = p.memberCID
	}

	client, err := falcon.NewClient(&apiConfig)
	if err != nil {
		return fmt.Errorf("failed to create CrowdStrike Falcon client: %w", err)
	}

	p.client = client
	return nil
}

// GetName returns the provider name
func (p *CrowdStrikeProvider) GetName() string {
	return "crowdstrike"
}

// GetSource returns the Terraform provider source
func (p *CrowdStrikeProvider) GetSource() string {
	return "crowdstrike/crowdstrike"
}

// GetConfig returns the provider configuration for generated Terraform files
func (p *CrowdStrikeProvider) GetConfig() cty.Value {
	config := map[string]cty.Value{
		"client_id":     cty.StringVal(p.clientID),
		"client_secret": cty.StringVal(p.clientSecret),
		"cloud":         cty.StringVal(p.cloud),
	}

	if p.memberCID != "" {
		config["member_cid"] = cty.StringVal(p.memberCID)
	}

	return cty.ObjectVal(config)
}

// InitService initializes a service with the shared client
func (p *CrowdStrikeProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"client":        p.client, // Pass the shared client instance
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"cloud":         p.cloud,
		"member_cid":    p.memberCID,
	})
	return nil
}

// GetSupportedService returns map of all supported services
func (p *CrowdStrikeProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		// Priority 1: Core Foundation
		"host_group": &HostGroupGenerator{},

		// Priority 2: Prevention Policies
		"prevention_policy_windows":         &PreventionPolicyWindowsGenerator{},
		"prevention_policy_linux":           &PreventionPolicyLinuxGenerator{},
		"prevention_policy_mac":             &PreventionPolicyMacGenerator{},
		"default_prevention_policy_windows": &DefaultPreventionPolicyWindowsGenerator{},
		"default_prevention_policy_linux":   &DefaultPreventionPolicyLinuxGenerator{},
		"default_prevention_policy_mac":     &DefaultPreventionPolicyMacGenerator{},
		"prevention_policy_precedence":      &PreventionPolicyPrecedenceGenerator{},
		"prevention_policy_attachment":      &PreventionPolicyAttachmentGenerator{},

		// Priority 3: Sensor Update Policies
		"sensor_update_policy":                       &SensorUpdatePolicyGenerator{},
		"default_sensor_update_policy":               &DefaultSensorUpdatePolicyGenerator{},
		"sensor_update_policy_precedence":            &SensorUpdatePolicyPrecedenceGenerator{},
		"sensor_update_policy_host_group_attachment": &SensorUpdatePolicyAttachmentGenerator{},

		// Priority 4: Cloud Security
		"cloud_aws_account":  &CloudAWSAccountGenerator{},
		"cloud_azure_tenant": &CloudAzureTenantGenerator{},

		// Priority 5: Content Update Policies
		"content_update_policy":            &ContentUpdatePolicyGenerator{},
		"default_content_update_policy":    &DefaultContentUpdatePolicyGenerator{},
		"content_update_policy_precedence": &ContentUpdatePolicyPrecedenceGenerator{},

		// Priority 6: File Integrity Monitoring
		"filevantage_policy":            &FilevantagePolicyGenerator{},
		"filevantage_rule_group":        &FilevantageRuleGroupGenerator{},
		"filevantage_policy_precedence": &FilevantagePolicyPrecedenceGenerator{},

		// Priority 7: Cloud Security Rules & Compliance
		"cloud_security_rule":               &CloudSecurityRuleGenerator{},
		"cloud_compliance_custom_framework": &CloudComplianceCustomFrameworkGenerator{},

		// Priority 8: Exclusions
		"sensor_visibility_exclusion": &SensorVisibilityExclusionGenerator{},

		// Priority 9: IT Automation
		"it_automation_policy":            &ItAutomationPolicyGenerator{},
		"it_automation_default_policy":    &ItAutomationDefaultPolicyGenerator{},
		"it_automation_policy_precedence": &ItAutomationPolicyPrecedenceGenerator{},
		"it_automation_task":              &ItAutomationTaskGenerator{},
		"it_automation_task_group":        &ItAutomationTaskGroupGenerator{},
	}
}

// GetResourceConnections returns resource relationship mappings for proper reference generation
func (p CrowdStrikeProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"prevention_policy_attachment": {
			"host_group": []string{"host_group_id", "id"},
		},
		"sensor_update_policy_host_group_attachment": {
			"host_group": []string{"host_group_id", "id"},
		},
		"content_update_policy": {
			"host_group": []string{"host_group_ids", "id"},
		},
		"filevantage_policy": {
			"host_group":             []string{"host_group_ids", "id"},
			"filevantage_rule_group": []string{"rule_group_ids", "id"},
		},
	}
}

// GetProviderData returns additional provider data if needed
func (p CrowdStrikeProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{}
}
