package secretsmanager

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/aws"

type AWSCryptoEngine struct {
	aws.AWSSDKConfig `mapstructure:",squash"`
	ID               string                 `mapstructure:"id"`
	Metadata         map[string]interface{} `mapstructure:"metadata"`
}
