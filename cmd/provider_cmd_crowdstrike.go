// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cmd

import (
	"errors"
	"os"

	crowdstrike_terraforming "github.com/GoogleCloudPlatform/terraformer/providers/crowdstrike"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/spf13/cobra"
)

func newCmdCrowdStrikeImporter(options ImportOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crowdstrike",
		Short: "Import current state to terraform configuration from CrowdStrike",
		Long:  "Import current state to terraform configuration from CrowdStrike",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientID := os.Getenv("FALCON_CLIENT_ID")
			if len(clientID) == 0 {
				return errors.New("FALCON_CLIENT_ID must be set")
			}
			clientSecret := os.Getenv("FALCON_CLIENT_SECRET")
			if len(clientSecret) == 0 {
				return errors.New("FALCON_CLIENT_SECRET must be set")
			}
			cloud := os.Getenv("FALCON_CLOUD")
			if len(cloud) == 0 {
				cloud = "autodiscover"
			}

			provider := newCrowdStrikeProvider()
			err := Import(provider, options, []string{clientID, clientSecret, cloud})
			if err != nil {
				return err
			}
			return nil
		},
	}
	cmd.AddCommand(listCmd(newCrowdStrikeProvider()))
	baseProviderFlags(cmd.PersistentFlags(), &options, "host_group", "host_group=group1:group2:group3")
	return cmd
}

func newCrowdStrikeProvider() terraformutils.ProviderGenerator {
	return &crowdstrike_terraforming.CrowdStrikeProvider{}
}
