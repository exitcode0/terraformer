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
	"regexp"
	"strings"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

// Helper functions for CrowdStrike provider

// createSimpleResource is a helper to create a resource with standard naming
func createSimpleResource(id, name, resourceType string, allowEmptyValues []string) terraformutils.Resource {
	// Sanitize the resource name to be Terraform-compatible
	sanitizedName := sanitizeResourceName(name)

	return terraformutils.NewSimpleResource(
		id,
		sanitizedName,
		resourceType,
		"crowdstrike",
		allowEmptyValues,
	)
}

// sanitizeResourceName converts a name into a Terraform-compatible resource name
// Replaces spaces, special characters, and ensures it starts with a letter
func sanitizeResourceName(name string) string {
	// Validate input
	if name == "" {
		return "unnamed_resource"
	}

	// Limit length to reasonable bounds
	if len(name) > 100 {
		name = name[:100]
	}

	// Replace spaces and hyphens with underscores
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Remove or replace other special characters
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)

	// Remove duplicate underscores
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")

	// Remove leading/trailing underscores
	name = strings.Trim(name, "_")

	// Ensure it starts with a letter (Terraform requirement)
	if len(name) > 0 && (name[0] >= '0' && name[0] <= '9') {
		name = "r_" + name
	}

	// Fallback if name becomes empty after sanitization
	if name == "" {
		name = "unnamed_resource"
	}

	// Convert to lowercase
	name = strings.ToLower(name)

	return name
}

// getFilteredIDs checks if specific IDs are requested via filters
// Returns (filteredIDs, shouldFilter)
// If shouldFilter is true, only the filteredIDs should be processed
// Validates filter values and returns only valid, non-empty IDs
func getFilteredIDs(generator *terraformutils.Service, resourceType string) ([]string, bool) {
	var filteredIDs []string

	for _, filter := range generator.Filter {
		if filter.FieldPath == "id" && filter.IsApplicable(resourceType) {
			for _, value := range filter.AcceptableValues {
				// Validate and sanitize filter value
				if validID := validateFilterID(value); validID != "" {
					filteredIDs = append(filteredIDs, validID)
				}
			}
		}
	}

	return filteredIDs, len(filteredIDs) > 0
}

// validateFilterID validates a filter ID value
// Returns empty string if invalid, sanitized ID if valid
func validateFilterID(id string) string {
	// Remove leading/trailing whitespace
	id = strings.TrimSpace(id)

	// Reject empty strings
	if id == "" {
		return ""
	}

	// Reject IDs that are too long (likely malformed)
	if len(id) > 200 {
		return ""
	}

	// Reject IDs with invalid characters (basic validation)
	// CrowdStrike IDs typically contain alphanumeric and hyphens
	for _, char := range id {
		if !isValidIDChar(char) {
			return ""
		}
	}

	return id
}

// isValidIDChar checks if a character is valid in a resource ID
func isValidIDChar(char rune) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '-' ||
		char == '_'
}

// createResourcesFromIDs creates resources from a list of IDs
func createResourcesFromIDs(ids []string, resourceType, prefix string, allowEmptyValues []string) []terraformutils.Resource {
	var resources []terraformutils.Resource

	for _, id := range ids {
		if id == "" {
			continue
		}

		// Create a unique resource name using the ID
		resourceName := fmt.Sprintf("%s_%s", prefix, sanitizeResourceName(id))

		resources = append(resources, terraformutils.NewSimpleResource(
			id,
			resourceName,
			resourceType,
			"crowdstrike",
			allowEmptyValues,
		))
	}

	return resources
}

// filterStringSlice filters a slice of strings to only include values in the filter list
func filterStringSlice(items []string, filterList []string) []string {
	if len(filterList) == 0 {
		return items
	}

	filterMap := make(map[string]bool)
	for _, f := range filterList {
		filterMap[f] = true
	}

	var filtered []string
	for _, item := range items {
		if filterMap[item] {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// createResourceWithName creates a resource with a custom name and ID
func createResourceWithName(id, name, resourceType string, allowEmptyValues []string) terraformutils.Resource {
	sanitizedName := sanitizeResourceName(name)

	return terraformutils.NewSimpleResource(
		id,
		sanitizedName,
		resourceType,
		"crowdstrike",
		allowEmptyValues,
	)
}

// nilStringValue safely gets the value of a string pointer, returning empty string if nil
func nilStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// stringPtr converts a string to a pointer
func stringPtr(s string) *string {
	return &s
}

// Common allow empty values for different resource types
var (
	// HostGroupAllowEmptyValues specifies fields that can be empty in host groups
	HostGroupAllowEmptyValues = []string{"tags."}

	// PreventionPolicyAllowEmptyValues specifies fields that can be empty in prevention policies
	PreventionPolicyAllowEmptyValues = []string{"tags.", "description"}

	// SensorUpdatePolicyAllowEmptyValues specifies fields that can be empty in sensor update policies
	SensorUpdatePolicyAllowEmptyValues = []string{"tags.", "description"}

	// ContentUpdatePolicyAllowEmptyValues specifies fields that can be empty in content update policies
	ContentUpdatePolicyAllowEmptyValues = []string{"tags.", "description"}

	// FilevantageAllowEmptyValues specifies fields that can be empty in FileVantage resources
	FilevantageAllowEmptyValues = []string{"tags.", "description"}

	// CloudAllowEmptyValues specifies fields that can be empty in cloud resources
	CloudAllowEmptyValues = []string{"tags."}

	// ITAutomationAllowEmptyValues specifies fields that can be empty in IT Automation resources
	ITAutomationAllowEmptyValues = []string{"tags.", "description"}
)
