# CrowdStrike

Use the `crowdstrike` provider to import existing CrowdStrike Falcon resources into Terraform.

## Requirements

* CrowdStrike Falcon API credentials (Client ID and Secret)
* Appropriate API scopes for the resources you want to import

## Authentication

The CrowdStrike provider requires authentication via API credentials. Set the following environment variables:

```sh
export FALCON_CLIENT_ID="your-client-id"
export FALCON_CLIENT_SECRET="your-client-secret"
export FALCON_CLOUD="us-1"  # Optional: us-1, us-2, eu-1, us-gov-1, or autodiscover (default)
```

For MSSP scenarios, you can optionally set:

```sh
export FALCON_MEMBER_CID="member-customer-id"
```

## Supported Resources

### Core Foundation
* `host_group` - Host groups for organizing endpoints

### Security Policies
* `prevention_policy_windows` - Windows prevention policies
* `prevention_policy_linux` - Linux prevention policies
* `prevention_policy_mac` - Mac prevention policies
* `default_prevention_policy_windows` - Default Windows prevention policy
* `default_prevention_policy_linux` - Default Linux prevention policy
* `default_prevention_policy_mac` - Default Mac prevention policy
* `prevention_policy_precedence` - Prevention policy precedence order
* `prevention_policy_attachment` - Prevention policy to host group attachments

### Sensor Management
* `sensor_update_policy` - Sensor update policies
* `default_sensor_update_policy` - Default sensor update policy
* `sensor_update_policy_precedence` - Sensor update policy precedence order
* `sensor_update_policy_host_group_attachment` - Sensor update policy to host group attachments

### Cloud Security
* `cloud_aws_account` - AWS account registration for cloud security
* `cloud_azure_tenant` - Azure tenant registration for cloud security

### Content Management
* `content_update_policy` - Content update policies
* `default_content_update_policy` - Default content update policy
* `content_update_policy_precedence` - Content update policy precedence order

### File Integrity Monitoring
* `filevantage_policy` - FileVantage FIM policies
* `filevantage_rule_group` - FileVantage rule groups
* `filevantage_policy_precedence` - FileVantage policy precedence order

### Cloud Security Rules & Compliance
* `cloud_security_rule` - Cloud security custom rules
* `cloud_compliance_custom_framework` - Cloud compliance custom frameworks

### Exclusions
* `sensor_visibility_exclusion` - Sensor visibility exclusions

### IT Automation
* `it_automation_policy` - IT automation policies
* `it_automation_default_policy` - Default IT automation policy
* `it_automation_policy_precedence` - IT automation policy precedence order
* `it_automation_task` - IT automation tasks
* `it_automation_task_group` - IT automation task groups

## Example Usage

### Import All Resources

```sh
terraformer import crowdstrike --resources="*" --path-output=./generated
```

### Import Specific Resources

```sh
# Import host groups
terraformer import crowdstrike --resources=host_group

# Import prevention policies
terraformer import crowdstrike --resources=prevention_policy_windows,prevention_policy_linux

# Import sensor update policies
terraformer import crowdstrike --resources=sensor_update_policy

# Import cloud security resources
terraformer import crowdstrike --resources=cloud_aws_account,cloud_azure_tenant
```

### Import with Filtering

```sh
# Import specific host groups by ID
terraformer import crowdstrike \
  --resources=host_group \
  --filter="host_group=id1:id2:id3"
```

### List Available Resources

```sh
terraformer import crowdstrike list
```

## API Scopes Required

The CrowdStrike API credentials must have appropriate scopes for the resources you want to import:

### Host Groups
- Host groups: Read

### Prevention Policies
- Prevention policies: Read

### Sensor Update Policies
- Sensor update policies: Read

### Cloud Security
- CSPM registration: Read
- Cloud security: Read

### Content Update Policies
- Sensor update policies: Read (shared with sensor updates)

### FileVantage
- Falcon FileVantage: Read

### Cloud Security Rules
- Custom IOA rules: Read

### Cloud Compliance
- Cloud Security Policies: Read

### Sensor Visibility Exclusions
- Sensor Visibility Exclusions: Read

### IT Automation
- IT Automation - Policies: Read
- IT Automation - Tasks: Read
- IT Automation - User Groups: Read

## Notes

* The provider will automatically handle pagination for large result sets
* Default policies (platform_default) are imported separately from custom policies
* Precedence resources manage the ordering of policies and require manual configuration
* Some resources like policy attachments create individual resources for each attachment relationship

## Output Structure

By default, Terraformer generates files in the following structure:

```
generated/
└── crowdstrike/
    ├── host_group/
    │   ├── host_group.tf
    │   └── terraform.tfstate
    ├── prevention_policy_windows/
    │   ├── prevention_policy_windows.tf
    │   └── terraform.tfstate
    └── ...
```

You can customize this structure with the `--path-pattern` flag.

## Troubleshooting

### Authentication Errors

If you see authentication errors, verify:
1. Your `FALCON_CLIENT_ID` and `FALCON_CLIENT_SECRET` are correct
2. Your API credentials have not expired
3. The `FALCON_CLOUD` setting matches your CrowdStrike instance

### Missing Resources

If expected resources are not imported:
1. Verify your API credentials have the required scopes
2. Check that the resources exist in your CrowdStrike tenant
3. Review the Terraformer output for any error messages

### Scope Errors

If you receive 403 Forbidden errors:
1. Review the API scopes listed above
2. Update your API credentials to include the required scopes
3. Create new API credentials if needed with appropriate permissions

## Further Information

For more details on the CrowdStrike Terraform Provider, see:
* [CrowdStrike Terraform Provider Documentation](https://registry.terraform.io/providers/CrowdStrike/crowdstrike/latest/docs)
* [CrowdStrike Falcon API Documentation](https://falcon.crowdstrike.com/documentation/page/a2a7fc0e/crowdstrike-oauth2-based-apis)

