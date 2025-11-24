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

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/crowdstrike/gofalcon/falcon/client"
)

type CrowdStrikeService struct { //nolint
	terraformutils.Service
}

// Client returns the shared CrowdStrike API client instance
// This is more efficient than creating a new client for each service
func (s *CrowdStrikeService) Client() (context.Context, *client.CrowdStrikeAPISpecification, error) {
	// Retrieve the pre-initialized client from Args
	client, ok := s.Args["client"].(*client.CrowdStrikeAPISpecification)
	if !ok || client == nil {
		// Fallback: this shouldn't happen if provider is correctly initialized
		return nil, nil, terraformutils.NewError("CrowdStrike client not properly initialized")
	}

	ctx := context.Background()
	return ctx, client, nil
}
