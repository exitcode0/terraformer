# CrowdStrike Terraformer - Quick Reference

Fast reference guide for common commands and patterns.

---

## 🚀 Quick Start

```bash
# Set credentials
export FALCON_CLIENT_ID="your-client-id"
export FALCON_CLIENT_SECRET="your-client-secret"
export FALCON_CLOUD="us-1"  # Optional

# Import host groups
terraformer import crowdstrike --resources="host_group"
```

---

## 📋 Common Commands

### List Available Resources
```bash
terraformer import crowdstrike list
```

### Import Single Resource Type
```bash
terraformer import crowdstrike --resources="host_group"
terraformer import crowdstrike --resources="prevention_policy_windows"
terraformer import crowdstrike --resources="sensor_update_policy"
```

### Import Multiple Resources
```bash
terraformer import crowdstrike --resources="host_group,prevention_policy_windows,sensor_update_policy"
```

### Import All Resources
```bash
terraformer import crowdstrike --resources="*"
```

### Import with Custom Output Directory
```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --path-output="./my-crowdstrike-terraform"
```

### Import with Verbose Output
```bash
terraformer import crowdstrike \
  --resources="host_group" \
  --verbose
```

---

## 🔐 Authentication Options

### Environment Variables (Recommended for CI/CD)
```bash
export FALCON_CLIENT_ID="your-client-id"
export FALCON_CLIENT_SECRET="your-client-secret"
export FALCON_CLOUD="us-1"
export FALCON_MEMBER_CID="member-cid"  # Optional for MSSP

terraformer import crowdstrike --resources="host_group"
```

### CLI Arguments (Recommended for Local Use)
```bash
terraformer import crowdstrike \
  --client-id="your-client-id" \
  --client-secret="your-client-secret" \
  --cloud="us-1" \
  --resources="host_group"
```

### Cloud Options
| Value | Description |
|-------|-------------|
| `us-1` | US Commercial 1 |
| `us-2` | US Commercial 2 |
| `eu-1` | EU Commercial |
| `us-gov-1` | US GovCloud |
| `autodiscover` | Auto-detect (default) |

---

## 🎯 Filtering (NEW!)

### Filter by Resource ID
```bash
# Single resource type with specific IDs
terraformer import crowdstrike \
  --resources="host_group" \
  --filter="host_group=hg-abc123:hg-def456:hg-xyz789"
```

### Filter Multiple Resource Types
```bash
# Multiple resource types, each with filters
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_windows" \
  --filter="host_group=hg-123;prevention_policy_windows=pp-456"
```

### Supported Filters

**✅ All major resources now support filtering!**

**Core:**
- `host_group=id1:id2`

**Prevention Policies:**
- `prevention_policy_windows=id1:id2`
- `prevention_policy_linux=id1:id2`
- `prevention_policy_mac=id1:id2`
- `prevention_policy_precedence=windows:linux:mac`
- `prevention_policy_attachment=policy_id1:policy_id2`

**Sensor Update:**
- `sensor_update_policy=id1:id2`
- `sensor_update_policy_precedence=windows:linux`
- `sensor_update_policy_host_group_attachment=policy_id1`

**Content Update:**
- `content_update_policy=id1:id2`
- `content_update_policy_precedence=windows:linux:mac`

**Exclusions:**
- `sensor_visibility_exclusion=id1:id2`

---

## 📚 Resource Categories

### Core Resources
```bash
terraformer import crowdstrike --resources="host_group"
```

### Prevention Policies
```bash
# All platforms
terraformer import crowdstrike --resources="prevention_policy_windows,prevention_policy_linux,prevention_policy_mac"

# With defaults
terraformer import crowdstrike --resources="prevention_policy_windows,default_prevention_policy_windows"

# With precedence and attachments
terraformer import crowdstrike --resources="prevention_policy_windows,prevention_policy_precedence,prevention_policy_attachment"
```

### Sensor Update Policies
```bash
# Basic
terraformer import crowdstrike --resources="sensor_update_policy"

# With defaults
terraformer import crowdstrike --resources="sensor_update_policy,default_sensor_update_policy"

# With precedence and attachments
terraformer import crowdstrike --resources="sensor_update_policy,sensor_update_policy_precedence,sensor_update_policy_host_group_attachment"
```

### Content Update Policies
```bash
terraformer import crowdstrike --resources="content_update_policy,default_content_update_policy,content_update_policy_precedence"
```

### Cloud Security
```bash
# AWS and Azure
terraformer import crowdstrike --resources="cloud_aws_account,cloud_azure_tenant"

# With security rules
terraformer import crowdstrike --resources="cloud_aws_account,cloud_security_rule"

# With compliance
terraformer import crowdstrike --resources="cloud_aws_account,cloud_compliance_custom_framework"
```

### File Integrity Monitoring
```bash
terraformer import crowdstrike --resources="filevantage_policy,filevantage_rule_group,filevantage_policy_precedence"
```

### IT Automation
```bash
terraformer import crowdstrike --resources="it_automation_task,it_automation_task_group,it_automation_policy"
```

### Exclusions
```bash
terraformer import crowdstrike --resources="sensor_visibility_exclusion"
```

---

## 🔄 After Import

### Navigate to Output
```bash
cd generated/crowdstrike
```

### Review Generated Files
```bash
ls -la
cat host_group.tf
cat provider.tf
```

### Initialize Terraform
```bash
terraform init
```

### Validate Configuration
```bash
terraform validate
```

### Plan Changes
```bash
terraform plan
```

### Apply (if needed)
```bash
terraform apply
```

---

## 🎨 Common Patterns

### Import Full Security Stack
```bash
terraformer import crowdstrike --resources="\
host_group,\
prevention_policy_windows,\
prevention_policy_linux,\
prevention_policy_mac,\
sensor_update_policy,\
content_update_policy"
```

### Import with Dependencies
```bash
# Host groups first, then policies, then attachments
terraformer import crowdstrike --resources="host_group"
terraformer import crowdstrike --resources="prevention_policy_windows,prevention_policy_linux"
terraformer import crowdstrike --resources="prevention_policy_attachment"
```

### Import for Specific Environment
```bash
# Use filters to import only production resources
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_windows" \
  --filter="host_group=prod-hg-1:prod-hg-2;prevention_policy_windows=prod-pp-1"
```

---

## 🐛 Troubleshooting

### Authentication Errors
```bash
# Verify credentials are set
echo $FALCON_CLIENT_ID
echo $FALCON_CLIENT_SECRET

# Try with verbose logging
terraformer import crowdstrike --resources="host_group" --verbose
```

### Empty Results
```bash
# Check without filters first
terraformer import crowdstrike --resources="host_group" --verbose

# Verify API scopes in CrowdStrike console
```

### Permission Errors
```bash
# Check API client has required scopes
# See README.md for required scopes per resource
```

---

## 📊 Examples by Use Case

### Audit Existing Configuration
```bash
# Import everything to review current state
terraformer import crowdstrike --resources="*" --path-output="./audit"
cd audit/crowdstrike
terraform plan  # Shows what would change
```

### Migrate to Terraform
```bash
# 1. Import resources
terraformer import crowdstrike --resources="host_group,prevention_policy_windows"

# 2. Review generated files
cd generated/crowdstrike
cat *.tf

# 3. Customize as needed
# Edit .tf files

# 4. Initialize and plan
terraform init
terraform plan
```

### Backup Configuration
```bash
# Regular backups of CrowdStrike config
DATE=$(date +%Y%m%d)
terraformer import crowdstrike \
  --resources="*" \
  --path-output="./backups/$DATE"
```

### Selective Import for Testing
```bash
# Import only test resources for validation
terraformer import crowdstrike \
  --resources="host_group,prevention_policy_windows" \
  --filter="host_group=test-hg-1;prevention_policy_windows=test-pp-1"
```

---

## 🔗 Quick Links

- **Full Documentation**: [README.md](README.md)
- **Examples**: Run `./example.sh`
- **Changelog**: [CHANGELOG.md](CHANGELOG.md)
- **Improvements**: [IMPROVEMENTS.md](IMPROVEMENTS.md)

---

## 💡 Pro Tips

1. **Start Small**: Import one resource type first to verify setup
2. **Use Filters**: In large environments, use filters to import specific resources
3. **Version Control**: Commit generated files to track changes over time
4. **Separate Imports**: Import different environments to different directories
5. **Review First**: Always run `terraform plan` before `terraform apply`
6. **Backup**: Keep backups of generated state files
7. **Incremental**: Import resources incrementally by priority
8. **Documentation**: Check API scopes if imports are empty

---

## 🚨 Common Mistakes to Avoid

❌ **Don't** import everything at once in a large environment  
✅ **Do** use filters to import incrementally

❌ **Don't** forget to set `FALCON_CLOUD` for non-US environments  
✅ **Do** set the correct cloud region

❌ **Don't** run `terraform apply` without reviewing `terraform plan`  
✅ **Do** review all changes before applying

❌ **Don't** mix production and test resources in the same import  
✅ **Do** use separate output directories for different environments

---

## 🎓 Learning Path

1. **Day 1**: Import host groups
   ```bash
   terraformer import crowdstrike --resources="host_group"
   ```

2. **Day 2**: Import prevention policies
   ```bash
   terraformer import crowdstrike --resources="prevention_policy_windows"
   ```

3. **Day 3**: Import with filters
   ```bash
   terraformer import crowdstrike --resources="host_group" --filter="host_group=id1"
   ```

4. **Day 4**: Import full security stack
   ```bash
   terraformer import crowdstrike --resources="host_group,prevention_policy_windows,sensor_update_policy"
   ```

5. **Day 5**: Use CLI arguments
   ```bash
   terraformer import crowdstrike --client-id="..." --client-secret="..." --resources="host_group"
   ```

---

**Need more details?** See [README.md](README.md) for comprehensive documentation.

