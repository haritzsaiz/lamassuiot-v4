package aws

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"

type AWSAuthenticationMethod string

const (
	Static     AWSAuthenticationMethod = "static"
	Default    AWSAuthenticationMethod = "default"
	AssumeRole AWSAuthenticationMethod = "role"
)

type AWSSDKConfig struct {
	AWSAuthenticationMethod AWSAuthenticationMethod `mapstructure:"auth_method"`
	EndpointURL             string                  `mapstructure:"endpoint_url"`
	AccessKeyID             string                  `mapstructure:"access_key_id"`
	SecretAccessKey         config.Password         `mapstructure:"secret_access_key"`
	SessionToken            config.Password         `mapstructure:"session_token"`
	Region                  string                  `mapstructure:"region"`
	RoleARN                 string                  `mapstructure:"role_arn"`
}
