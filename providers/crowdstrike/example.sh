#!/bin/bash

# CrowdStrike Terraformer Example Script
# This script demonstrates various ways to use the CrowdStrike provider with Terraformer
# 
# Usage: ./example.sh
#
# Prerequisites:
# - Set FALCON_CLIENT_ID environment variable
# - Set FALCON_CLIENT_SECRET environment variable
# - Optionally set FALCON_CLOUD (defaults to autodiscover)

set -e

echo "========================================="
echo " CrowdStrike Terraformer Example Script"
echo "========================================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if required environment variables are set
if [ -z "$FALCON_CLIENT_ID" ] || [ -z "$FALCON_CLIENT_SECRET" ]; then
    echo -e "${RED}Error: Required environment variables not set${NC}"
    echo ""
    echo "Please set the following environment variables:"
    echo ""
    echo "  export FALCON_CLIENT_ID=\"your-client-id\""
    echo "  export FALCON_CLIENT_SECRET=\"your-client-secret\""
    echo "  export FALCON_CLOUD=\"us-1\"  # Optional: us-1, us-2, eu-1, us-gov-1, or autodiscover"
    echo "  export FALCON_MEMBER_CID=\"member-cid\"  # Optional: For MSSP scenarios"
    echo ""
    exit 1
fi

# Set default cloud if not specified
if [ -z "$FALCON_CLOUD" ]; then
    export FALCON_CLOUD="autodiscover"
    echo -e "${YELLOW}No FALCON_CLOUD specified, using: autodiscover${NC}"
else
    echo -e "${GREEN}Using CrowdStrike Cloud: $FALCON_CLOUD${NC}"
fi

echo ""

# Create output directory
OUTPUT_DIR="./crowdstrike-terraform-import"
mkdir -p "$OUTPUT_DIR"

echo -e "${BLUE}Output directory: $OUTPUT_DIR${NC}"
echo ""

# Function to run terraformer with error handling
run_import() {
    local resource=$1
    local description=$2
    
    echo "----------------------------------------"
    echo -e "${BLUE}$description${NC}"
    echo "----------------------------------------"
    
    if terraformer import crowdstrike \
        --resources="$resource" \
        --path-output="$OUTPUT_DIR" \
        --verbose; then
        echo -e "${GREEN}✓ Successfully imported $resource${NC}"
    else
        echo -e "${RED}✗ Failed to import $resource${NC}"
        return 1
    fi
    echo ""
}

# Example 1: List available resources
echo "========================================="
echo " Example 1: List Available Resources"
echo "========================================="
echo ""
terraformer import crowdstrike list
echo ""

# Example 2: Import Host Groups
echo "========================================="
echo " Example 2: Import Host Groups"
echo "========================================="
echo ""
run_import "host_group" "Importing Host Groups (foundational resource)"

# Example 3: Import Prevention Policies
echo "========================================="
echo " Example 3: Import Prevention Policies"
echo "========================================="
echo ""
echo "Importing prevention policies for all platforms..."
echo ""
run_import "prevention_policy_windows" "Importing Windows Prevention Policies"
run_import "prevention_policy_linux" "Importing Linux Prevention Policies"
run_import "prevention_policy_mac" "Importing macOS Prevention Policies"

# Example 4: Import Default Policies
echo "========================================="
echo " Example 4: Import Default Policies"
echo "========================================="
echo ""
run_import "default_prevention_policy_windows" "Importing Default Windows Prevention Policy"
run_import "default_sensor_update_policy" "Importing Default Sensor Update Policy"

# Example 5: Import Policy Attachments
echo "========================================="
echo " Example 5: Import Policy Attachments"
echo "========================================="
echo ""
run_import "prevention_policy_attachment" "Importing Prevention Policy-Host Group Attachments"

# Example 6: Import Sensor Update Policies
echo "========================================="
echo " Example 6: Import Sensor Update Policies"
echo "========================================="
echo ""
run_import "sensor_update_policy" "Importing Sensor Update Policies"
run_import "sensor_update_policy_host_group_attachment" "Importing Sensor Update Policy Attachments"

# Example 7: Import Cloud Resources (if available)
echo "========================================="
echo " Example 7: Import Cloud Resources"
echo "========================================="
echo ""
echo "Note: These may be empty if you haven't configured cloud integrations"
echo ""
run_import "cloud_aws_account" "Importing AWS Account Integrations" || echo -e "${YELLOW}No AWS accounts found${NC}"
run_import "cloud_azure_tenant" "Importing Azure Tenant Integrations" || echo -e "${YELLOW}No Azure tenants found${NC}"

# Example 8: Show generated files
echo "========================================="
echo " Example 8: Generated Files"
echo "========================================="
echo ""
echo "Terraform files generated in: $OUTPUT_DIR"
echo ""

if [ -d "$OUTPUT_DIR/crowdstrike" ]; then
    echo "Directory structure:"
    tree "$OUTPUT_DIR/crowdstrike" -L 2 2>/dev/null || find "$OUTPUT_DIR/crowdstrike" -maxdepth 2 -type f
    echo ""
    
    echo "File count by type:"
    echo "  .tf files:      $(find "$OUTPUT_DIR/crowdstrike" -name "*.tf" | wc -l)"
    echo "  .tfstate files: $(find "$OUTPUT_DIR/crowdstrike" -name "*.tfstate" | wc -l)"
    echo ""
else
    echo -e "${YELLOW}No crowdstrike directory created${NC}"
fi

# Example 9: Validate Generated Terraform
echo "========================================="
echo " Example 9: Validate Generated Terraform"
echo "========================================="
echo ""

if [ -d "$OUTPUT_DIR/crowdstrike" ]; then
    cd "$OUTPUT_DIR/crowdstrike"
    
    echo "Initializing Terraform..."
    if terraform init > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Terraform initialized${NC}"
        echo ""
        
        echo "Validating configuration..."
        if terraform validate; then
            echo -e "${GREEN}✓ Terraform configuration is valid${NC}"
        else
            echo -e "${RED}✗ Terraform validation failed${NC}"
        fi
    else
        echo -e "${YELLOW}Note: Terraform init failed - you may need to configure the provider${NC}"
    fi
    
    cd - > /dev/null
else
    echo -e "${YELLOW}No files to validate${NC}"
fi

echo ""

# Example 10: Show sample commands for next steps
echo "========================================="
echo " Next Steps"
echo "========================================="
echo ""
echo "Your CrowdStrike resources have been imported!"
echo ""
echo "To review the generated configuration:"
echo "  cd $OUTPUT_DIR/crowdstrike"
echo "  cat *.tf"
echo ""
echo "To plan changes:"
echo "  cd $OUTPUT_DIR/crowdstrike"
echo "  terraform init"
echo "  terraform plan"
echo ""
echo "To import specific resources with filters:"
echo "  terraformer import crowdstrike \\"
echo "    --resources=\"host_group\" \\"
echo "    --filter=\"host_group=id1:id2:id3\""
echo ""
echo "To import with CLI arguments instead of env vars:"
echo "  terraformer import crowdstrike \\"
echo "    --client-id=\"\$FALCON_CLIENT_ID\" \\"
echo "    --client-secret=\"\$FALCON_CLIENT_SECRET\" \\"
echo "    --cloud=\"us-1\" \\"
echo "    --resources=\"host_group\""
echo ""
echo "For more examples, see providers/crowdstrike/README.md"
echo ""

echo "========================================="
echo " Example Script Complete!"
echo "========================================="
echo ""
echo -e "${GREEN}All examples completed successfully!${NC}"
echo ""
echo "Generated files are in: $OUTPUT_DIR"
echo ""

