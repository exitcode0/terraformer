# CrowdStrike Provider for Terraformer

This provider allows you to import existing CrowdStrike Falcon resources into Terraform configuration files and state.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Supported Resources](#supported-resources)
- [Authentication](#authentication)
- [Usage Examples](#usage-examples)
- [Filtering](#filtering)
- [Resource Dependencies](#resource-dependencies)
- [Output Structure](#output-structure)
- [Troubleshooting](#troubleshooting)

## Prerequisites

### 1. CrowdStrike Falcon API Credentials

You need a CrowdStrike Falcon API client with appropriate scopes:

| Resource Type | Required API Scope |
|--------------|-------------------|
| Host Groups | `Host groups: Read` |
| Prevention Policies | `Prevention Policies: Read` |
| Sensor Update Policies | `Sensor update policies: Read` |
| Content Update Policies | `Content update policies: Read` |
| Cloud Accounts (AWS/Azure) | `CSPM registration: Read` |
| File Integrity Monitoring | `FileVantage: Read` |
| IT Automation | `IT Automation: Read` |
| Sensor Visibility Exclusions | `Sensor visibility exclusions: Read` |
| Cloud Security | `Cloud Security: Read` |
| Cloud Compliance | `Cloud Compliance: Read` |

### 2. Terraform CrowdStrike Provider

Ensure you have the CrowdStrike Terraform provider available:

```bash
terraform {
  required_providers {
    crowdstrike = {
      source  = "crowdstrike/crowdstrike"
      version = "~> 0.3"
    }
  }
}
```

## Supported Resources

### Core Resources
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `host_group` | `crowdstrike_host_group` | Host Groups for organizing endpoints |

### Prevention Policies
| Terraformer Resource | Terraform Resource Type | Platform |
|---------------------|------------------------|----------|
| `prevention_policy_windows` | `crowdstrike_prevention_policy_windows` | Windows |
| `prevention_policy_linux` | `crowdstrike_prevention_policy_linux` | Linux |
| `prevention_policy_mac` | `crowdstrike_prevention_policy_mac` | macOS |
| `default_prevention_policy_windows` | `crowdstrike_default_prevention_policy_windows` | Windows Default |
| `default_prevention_policy_linux` | `crowdstrike_default_prevention_policy_linux` | Linux Default |
| `default_prevention_policy_mac` | `crowdstrike_default_prevention_policy_mac` | macOS Default |
| `prevention_policy_precedence` | `crowdstrike_prevention_policy_precedence` | Policy Precedence |
| `prevention_policy_attachment` | `crowdstrike_prevention_policy_attachment` | Policy-Host Group Attachments |

### Sensor Update Policies
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `sensor_update_policy` | `crowdstrike_sensor_update_policy` | Sensor Update Policies |
| `default_sensor_update_policy` | `crowdstrike_default_sensor_update_policy` | Default Sensor Update Policy |
| `sensor_update_policy_precedence` | `crowdstrike_sensor_update_policy_precedence` | Policy Precedence |
| `sensor_update_policy_host_group_attachment` | `crowdstrike_sensor_update_policy_host_group_attachment` | Policy-Host Group Attachments |

### Content Update Policies
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `content_update_policy` | `crowdstrike_content_update_policy` | Content Update Policies |
| `default_content_update_policy` | `crowdstrike_default_content_update_policy` | Default Content Update Policy |
| `content_update_policy_precedence` | `crowdstrike_content_update_policy_precedence` | Policy Precedence |

### Cloud Security
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `cloud_aws_account` | `crowdstrike_cloud_aws_account` | AWS Account Integrations |
| `cloud_azure_tenant` | `crowdstrike_cloud_azure_tenant` | Azure Tenant Integrations |
| `cloud_security_rule` | `crowdstrike_cloud_security_custom_rule` | Custom Cloud Security Rules |
| `cloud_compliance_custom_framework` | `crowdstrike_cloud_compliance_custom_framework` | Custom Compliance Frameworks |

### File Integrity Monitoring (FIM)
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `filevantage_policy` | `crowdstrike_filevantage_policy` | FileVantage Policies |
| `filevantage_rule_group` | `crowdstrike_filevantage_rule_group` | FileVantage Rule Groups |
| `filevantage_policy_precedence` | `crowdstrike_filevantage_policy_precedence` | Policy Precedence |

### IT Automation
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `it_automation_policy` | `crowdstrike_it_automation_policy` | IT Automation Policies |
| `it_automation_default_policy` | `crowdstrike_it_automation_default_policy` | Default IT Automation Policy |
| `it_automation_policy_precedence` | `crowdstrike_it_automation_policy_precedence` | Policy Precedence |
| `it_automation_task` | `crowdstrike_it_automation_task` | IT Automation Tasks |
| `it_automation_task_group` | `crowdstrike_it_automation_task_group` | IT Automation Task Groups |

### Exclusions
| Terraformer Resource | Terraform Resource Type | Description |
|---------------------|------------------------|-------------|
| `sensor_visibility_exclusion` | `crowdstrike_sensor_visibility_exclusion` | Sensor Visibility Exclusions |

## Authentication

### Option 1: Environment Variables (Recommended for CI/CD)

```bash
export FALCON_CLIENT_ID="your-client-id"
export FALCON_CLIENT_SECRET="your-client-secret"
export FALCON_CLOUD="us-1"  # Optional: us-1, us-2, eu-1, us-gov-1, or autodiscover (default)
export FALCON_MEMBER_CID="member-cid"  # Optional: For MSSP scenarios
```

### Option 2: Command-line Flags

```bash
terraformer import crowdstrike \
  --client-id="your-client-id" \
  --client-secret="your-client-secret" \
  --cloud="us-1" \
  --resources="host_group"
```

### Cloud Region Options

| Value | Description |
|-------|-------------|
| `us-1` | US Commercial 1 |
| `us-2` | US Commercial 2 |
| `eu-1` | EU Commercial |
| `us-gov-1` | US GovCloud |
| `autodiscover` | Automatically detect (default) |

## Usage Examples

### Import All Supported Resources

```bash
terraformer import crowdstrike --resources="*"
```

### Import Specific Resource Types

```bash
# Import host groups only
terraformer import crowdstrike --resources="host_group"

# Import multiple resource types
terraformer import crowdstrike --resources="host_group,prevention_policy_windows,sensor_update_policy"
```

### Import with Custom Output Directory

```bash
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_windows" \
  --path-output="./crowdstrike-terraform"
```

### Import Prevention Policies by Platform

```bash
# Import only Windows prevention policies
terraformer import crowdstrike --resources="prevention_policy_windows,default_prevention_policy_windows"

# Import all Linux-related resources
terraformer import crowdstrike --resources="prevention_policy_linux,default_prevention_policy_linux"
```

### Import Cloud Security Resources

```bash
# Import AWS account integrations
terraformer import crowdstrike --resources="cloud_aws_account"

# Import Azure tenant integrations
terraformer import crowdstrike --resources="cloud_azure_tenant"

# Import all cloud resources
terraformer import crowdstrike --resources="cloud_aws_account,cloud_azure_tenant,cloud_security_rule"
```

### List Available Resources

```bash
terraformer import crowdstrike list
```

## Filtering

You can filter resources by ID to import only specific items:

### Filter by Resource ID

```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --filter="host_group=id1:id2:id3"
```

### Filter Multiple Resource Types

```bash
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_windows" \
  --filter="host_group=hg-id1:hg-id2;prevention_policy_windows=pp-id1:pp-id2"
```

### Common Use Cases

```bash
# Import specific host groups and their associated policies
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_attachment" \
  --filter="host_group=abc123:def456"

# Import only production sensor update policies
terraformer import crowdstrike \
  --resources="sensor_update_policy" \
  --filter="sensor_update_policy=prod-policy-id"

# Import precedence for specific platforms only
terraformer import crowdstrike \
  --resources="prevention_policy_precedence,sensor_update_policy_precedence" \
  --filter="prevention_policy_precedence=windows:linux;sensor_update_policy_precedence=windows"

# Import exclusions by ID
terraformer import crowdstrike \
  --resources="sensor_visibility_exclusion" \
  --filter="sensor_visibility_exclusion=exclusion-id-1:exclusion-id-2"
```

## Resource Dependencies

Terraformer automatically handles resource dependencies. When you import resources, it creates proper references between related resources:

### Automatic Dependency Resolution

```hcl
# Host Group
resource "crowdstrike_host_group" "example" {
  name = "Production Servers"
  # ... other attributes ...
}

# Prevention Policy with automatic reference
resource "crowdstrike_prevention_policy_attachment" "example" {
  host_group_id = crowdstrike_host_group.example.id  # ← Automatic reference
  policy_id     = "policy-id"
}
```

### Dependency Map

- **Prevention Policy Attachments** → Host Groups
- **Sensor Update Policy Attachments** → Host Groups
- **Content Update Policies** → Host Groups
- **FileVantage Policies** → Host Groups + Rule Groups
- **IT Automation Policies** → Host Groups + Task Groups

## Output Structure

After running terraformer, you'll get the following structure:

```
./generated/crowdstrike/
├── provider.tf                    # CrowdStrike provider configuration
├── host_group.tf                  # Host group resources
├── prevention_policy_windows.tf   # Windows prevention policies
├── prevention_policy_linux.tf     # Linux prevention policies
├── sensor_update_policy.tf        # Sensor update policies
├── outputs.tf                     # Outputs (optional)
└── terraform.tfstate              # Terraform state file
```

### Example Generated Provider Configuration

```hcl
provider "crowdstrike" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  cloud         = "us-1"
}
```

### Example Generated Resource

```hcl
resource "crowdstrike_host_group" "production_servers" {
  name               = "Production Servers"
  description        = "All production servers"
  group_type         = "static"
  assignment_rule    = ""
}
```

## Advanced Usage

### Import with Verbose Output

```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --verbose
```

### Import with Specific State File

```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --state="./terraform.tfstate"
```

### Import and Generate Plan

After importing:

```bash
cd generated/crowdstrike
terraform init
terraform plan
```

## Troubleshooting

### Authentication Errors

**Error**: `FALCON_CLIENT_ID is required`

**Solution**: Ensure credentials are set:
```bash
export FALCON_CLIENT_ID="your-client-id"
export FALCON_CLIENT_SECRET="your-client-secret"
```

### Insufficient API Scopes

**Error**: `403 Forbidden` or `access denied`

**Solution**: Verify your API client has the required scopes listed in [Prerequisites](#prerequisites).

### Empty Results

**Issue**: No resources are imported

**Possible Causes**:
1. No resources exist in your CrowdStrike tenant
2. API client lacks required scopes
3. Filters are too restrictive

**Solution**: Try importing without filters first:
```bash
terraformer import crowdstrike --resources="host_group" --verbose
```

### Cloud Region Issues

**Error**: `autodiscover failed`

**Solution**: Explicitly set the cloud region:
```bash
export FALCON_CLOUD="us-1"  # or us-2, eu-1, us-gov-1
```

### Rate Limiting

**Error**: `429 Too Many Requests`

**Solution**: The CrowdStrike API has rate limits. If you hit them:
1. Import resources in smaller batches
2. Add delays between imports
3. Contact CrowdStrike support for increased limits

## Best Practices

### 1. Start Small

Begin by importing a single resource type:
```bash
terraformer import crowdstrike --resources="host_group"
```

### 2. Use Filters in Large Environments

For tenants with many resources, use filters:
```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --filter="host_group=prod-hg1:prod-hg2"
```

### 3. Version Control

Always commit generated Terraform files to version control:
```bash
git add generated/crowdstrike/
git commit -m "Import CrowdStrike configuration"
```

### 4. Review Before Applying

Always review the generated configuration:
```bash
cd generated/crowdstrike
terraform init
terraform plan  # Review changes before applying
```

### 5. Incremental Imports

Import resources incrementally by priority:
1. First: Core resources (Host Groups)
2. Second: Policies (Prevention, Sensor Update)
3. Third: Attachments and relationships
4. Fourth: Advanced features (Cloud, FIM, IT Automation)

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Import CrowdStrike Configuration

on:
  workflow_dispatch:

jobs:
  import:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Terraformer
        run: |
          curl -LO https://github.com/GoogleCloudPlatform/terraformer/releases/download/$(curl -s https://api.github.com/repos/GoogleCloudPlatform/terraformer/releases/latest | grep tag_name | cut -d '"' -f 4)/terraformer-all-linux-amd64
          chmod +x terraformer-all-linux-amd64
          sudo mv terraformer-all-linux-amd64 /usr/local/bin/terraformer
      
      - name: Import CrowdStrike Resources
        env:
          FALCON_CLIENT_ID: ${{ secrets.FALCON_CLIENT_ID }}
          FALCON_CLIENT_SECRET: ${{ secrets.FALCON_CLIENT_SECRET }}
          FALCON_CLOUD: "us-1"
        run: |
          terraformer import crowdstrike --resources="host_group,prevention_policy_windows"
      
      - name: Commit Changes
        run: |
          git config --local user.email "action@github.com"
          git config --local user.name "GitHub Action"
          git add generated/
          git commit -m "Update CrowdStrike configuration" || echo "No changes"
          git push
```

## Additional Resources

- [CrowdStrike Terraform Provider Documentation](https://registry.terraform.io/providers/crowdstrike/crowdstrike/latest/docs)
- [CrowdStrike API Documentation](https://falcon.crowdstrike.com/documentation/page/a2a7fc0e/crowdstrike-api-specification-for-swagger)
- [Terraformer Documentation](https://github.com/GoogleCloudPlatform/terraformer)
- [Terraform Best Practices](https://www.terraform.io/docs/cloud/guides/recommended-practices/index.html)

## Contributing

Found a bug or want to contribute? Please open an issue or pull request in the [Terraformer repository](https://github.com/GoogleCloudPlatform/terraformer).

## License

This provider follows the same Apache 2.0 license as Terraformer.

