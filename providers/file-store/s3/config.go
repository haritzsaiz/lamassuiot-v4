package s3

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/aws"

type AWSS3FilesystemConfig struct {
	aws.AWSSDKConfig `mapstructure:",squash"`
	BucketName       string                 `mapstructure:"bucket_name"`
	ID               string                 `mapstructure:"id"`
	Metadata         map[string]interface{} `mapstructure:"metadata"`
}
