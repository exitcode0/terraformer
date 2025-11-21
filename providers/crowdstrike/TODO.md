# CrowdStrike Provider Implementation TODO

## Implementation Progress

### ✅ Foundation Setup
- [x] Create `crowdstrike_provider.go`
- [x] Create `crowdstrike_service.go`
- [x] Create `provider_cmd_crowdstrike.go`
- [x] Register in `cmd/root.go`
- [x] Update `go.mod` dependencies

### Priority 1: Core Foundation
- [ ] **Host Groups** (`host_group.go`)
  - Resource: `crowdstrike_host_group`

### Priority 2: Essential Security Policies
- [ ] **Prevention Policies** (`prevention_policy_*.go`)
  - Resource: `crowdstrike_prevention_policy_windows`
  - Resource: `crowdstrike_prevention_policy_linux`
  - Resource: `crowdstrike_prevention_policy_mac`
  - Resource: `crowdstrike_default_prevention_policy_windows`
  - Resource: `crowdstrike_default_prevention_policy_linux`
  - Resource: `crowdstrike_default_prevention_policy_mac`
  - Resource: `crowdstrike_prevention_policy_precedence`
  - Resource: `crowdstrike_prevention_policy_attachment`

### Priority 3: Sensor Management
- [ ] **Sensor Update Policies** (`sensor_update_policy_*.go`)
  - Resource: `crowdstrike_sensor_update_policy`
  - Resource: `crowdstrike_default_sensor_update_policy`
  - Resource: `crowdstrike_sensor_update_policy_precedence`
  - Resource: `crowdstrike_sensor_update_policy_host_group_attachment`
  - Data Source: `crowdstrike_sensor_update_policy_builds`

### Priority 4: Cloud Security (AWS/Azure)
- [ ] **Cloud AWS Account** (`cloud_aws_account.go`)
  - Resource: `crowdstrike_cloud_aws_account`
  - Data Source: `crowdstrike_cloud_aws_account` (data source)
- [ ] **Cloud Azure Tenant** (`cloud_azure_tenant.go`)
  - Resource: `crowdstrike_cloud_azure_tenant`
  - Resource: `crowdstrike_cloud_azure_tenant_eventhub_settings`

### Priority 5: Content Management
- [ ] **Content Update Policies** (`content_update_policy_*.go`)
  - Resource: `crowdstrike_content_update_policy`
  - Resource: `crowdstrike_default_content_update_policy`
  - Resource: `crowdstrike_content_update_policy_precedence`
  - Data Source: `crowdstrike_content_category_versions`

### Priority 6: File Integrity Monitoring
- [ ] **FileVantage** (`filevantage_*.go`)
  - Resource: `crowdstrike_filevantage_policy`
  - Resource: `crowdstrike_filevantage_rule_group`
  - Resource: `crowdstrike_filevantage_policy_precedence`

### Priority 7: Cloud Security Rules & Compliance
- [ ] **Cloud Security Rules** (`cloud_security_*.go`)
  - Resource: `crowdstrike_cloud_security_rule`
  - Data Source: `crowdstrike_cloud_security_rules`
- [ ] **Cloud Compliance** (`cloud_compliance_*.go`)
  - Resource: `crowdstrike_cloud_compliance_custom_framework`
  - Data Source: `crowdstrike_cloud_compliance_framework_controls`

### Priority 8: Exclusions
- [ ] **Sensor Visibility Exclusions** (`sensor_visibility_exclusion.go`)
  - Resource: `crowdstrike_sensor_visibility_exclusion`

### Priority 9: IT Automation
- [ ] **IT Automation** (`it_automation_*.go`)
  - Resource: `crowdstrike_it_automation_policy`
  - Resource: `crowdstrike_it_automation_default_policy`
  - Resource: `crowdstrike_it_automation_policy_precedence`
  - Resource: `crowdstrike_it_automation_task`
  - Resource: `crowdstrike_it_automation_task_group`

### Documentation
- [ ] Update `/docs/crowdstrike.md` with usage examples
- [ ] Update main `README.md` to list CrowdStrike provider

## Testing Checklist

After each priority level implementation:
1. Test basic import: `terraformer import crowdstrike --resources=<resource> --path-output=./test-output`
2. Verify `.tf` files are generated correctly
3. Verify `.tfstate` files are valid
4. Test with filters if applicable

## Notes

- All resources use the gofalcon SDK client
- Authentication via `FALCON_CLIENT_ID`, `FALCON_CLIENT_SECRET`, `FALCON_CLOUD`
- Optional `FALCON_MEMBER_CID` for MSSP scenarios
- Follow the pattern from Okta/Auth0 providers for consistency

