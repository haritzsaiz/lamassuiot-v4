package s3

import (
	"context"

	s3v2 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/aws"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"github.com/lamassuiot/lamassuiot/v4/providers/file-store/service"
	filestore "github.com/lamassuiot/lamassuiot/v4/providers/file-store/service"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
)

func Register() {
	filestore.RegisterProvider(service.AWSS3, func(logger *logger.Logger, conf service.FileStoreConfig) (*blob.Bucket, error) {
		engineConfig, _ := service.FileStoreConfigAdapter[AWSS3FilesystemConfig]{}.Marshal(conf)

		awsCfg, err := aws.GetAwsSdkConfig(engineConfig.Config.AWSSDKConfig)
		if err != nil {
			return nil, err
		}

		clientV2 := s3v2.NewFromConfig(*awsCfg)
		bucket, err := s3blob.OpenBucketV2(context.Background(), clientV2, engineConfig.Config.BucketName, nil)

		if err != nil {
			return nil, err
		}

		return bucket, nil
	})
}
