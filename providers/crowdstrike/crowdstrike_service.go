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
	"github.com/crowdstrike/gofalcon/falcon"
	"github.com/crowdstrike/gofalcon/falcon/client"
)

type CrowdStrikeService struct { //nolint
	terraformutils.Service
}

func (s *CrowdStrikeService) Client() (context.Context, *client.CrowdStrikeAPISpecification, error) {
	clientID := s.Args["client_id"].(string)
	clientSecret := s.Args["client_secret"].(string)
	cloud := s.Args["cloud"].(string)
	memberCID, _ := s.Args["member_cid"].(string)

	apiConfig := falcon.ApiConfig{
		ClientId:     clientID,
		ClientSecret: clientSecret,
		Cloud:        falcon.Cloud(cloud),
		Context:      context.Background(),
	}

	if memberCID != "" {
		apiConfig.MemberCID = memberCID
	}

	client, err := falcon.NewClient(&apiConfig)
	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()
	return ctx, client, nil
}
