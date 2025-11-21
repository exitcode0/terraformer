### Use with Okta

Example:

```
$ export OKTA_ORG_NAME=<ORG_NAME>
$ export OKTA_BASE_URL=<BASE_URL>
$ export OKTA_API_TOKEN=<API_TOKEN>
$ terraformer import okta --resources=okta_user,okta_group
```

If you login to your Okta instance at: https://dev-12345678.okta.com/ you would configure:
```
$ export OKTA_ORG_NAME=dev-12345678
$ export OKTA_BASE_URL=okta.com
```


List of supported Okta services:

*    `admin_role`
     * `okta_admin_role_custom`
     * `okta_admin_role_custom_assignments`
*    `app`
     * `okta_app_auto_login`
     * `okta_app_basic_auth`
     * `okta_app_bookmark`
     * `okta_app_group_assignments`
     * `okta_app_oauth`
     * `okta_app_oauth_api_scope`
     * `okta_app_oauth_role_assignment`
     * `okta_app_saml`
     * `okta_app_secure_password_store`
     * `okta_app_signon_policy`
     * `okta_app_signon_policy_rule`
     * `okta_app_swa`
     * `okta_app_three_field`
     * `okta_app_user_base_schema_property`
*    `authorization_server`
     * `okta_auth_server`
     * `okta_auth_server_claim`
     * `okta_auth_server_policy`
     * `okta_auth_server_policy_rule`
     * `okta_auth_server_scope`
*    `event_hook`
*    * `okta_event_hook`
*    `factor`
     * `okta_factor`
*    `group`
     * `okta_group`
     * `okta_group_memberships`
     * `okta_group_owner`
     * `okta_group_role`
     * `okta_group_rule`
     * `okta_group_schema_property`
*    `idp`
     * `okta_idp_oidc`
     * `okta_idp_saml`
     * `okta_idp_social`
*    `inline_hook`
*    * `okta_inline_hook`
*    `network_zone`
     * `okta_network_zone`
*    `policy`
     * `okta_policy_device_assurance_android`
     * `okta_policy_device_assurance_chromeos`
     * `okta_policy_device_assurance_ios`
     * `okta_policy_device_assurance_macos`
     * `okta_policy_device_assurance_windows`
     * `okta_policy_mfa`
     * `okta_policy_password`
     * `okta_policy_rule_mfa`
     * `okta_policy_rule_password`
     * `okta_policy_rule_signon`
     * `okta_policy_signon`
*    `profile_mapping`
     * `okta_profile_mapping`
*    `request_condition`
     * `okta_request_condition`
*    `resource_set`
     * `okta_resource_set`
*    `template_sms`
     * `okta_template_sms`
*    `trusted_origin`
     * `okta_trusted_origin`
*    `user`
     * `okta_user`
*    `user_type`
     * `okta_user_type`
