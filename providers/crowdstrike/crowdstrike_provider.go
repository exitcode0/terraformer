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
	"errors"
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/zclconf/go-cty/cty"
)

type CrowdStrikeProvider struct { //nolint
	terraformutils.Provider
	clientID     string
	clientSecret string
	cloud        string
	memberCID    string
}

func (p *CrowdStrikeProvider) Init(args []string) error {
	clientID := os.Getenv("FALCON_CLIENT_ID")
	if clientID == "" {
		return errors.New("set FALCON_CLIENT_ID env var")
	}
	p.clientID = clientID

	clientSecret := os.Getenv("FALCON_CLIENT_SECRET")
	if clientSecret == "" {
		return errors.New("set FALCON_CLIENT_SECRET env var")
	}
	p.clientSecret = clientSecret

	cloud := os.Getenv("FALCON_CLOUD")
	if cloud == "" {
		cloud = "autodiscover"
	}
	p.cloud = cloud

	// Member CID is optional for MSSP scenarios
	memberCID := os.Getenv("FALCON_MEMBER_CID")
	p.memberCID = memberCID

	return nil
}

func (p *CrowdStrikeProvider) GetName() string {
	return "crowdstrike"
}

func (p *CrowdStrikeProvider) GetSource() string {
	return "crowdstrike/crowdstrike"
}

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
		"client_id":     p.clientID,
		"client_secret": p.clientSecret,
		"cloud":         p.cloud,
		"member_cid":    p.memberCID,
	})
	return nil
}

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

func (p CrowdStrikeProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"prevention_policy_attachment": {
			"host_group": []string{"host_group_id", "id"},
		},
		"sensor_update_policy_host_group_attachment": {
			"host_group": []string{"host_group_id", "id"},
		},
	}
}

func (p CrowdStrikeProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{}
}
