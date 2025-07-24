package secretsmanager

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms/cryptoengines"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/aws"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

func Register() {
	cryptoengines.RegisterProvider(cryptoengines.AWSKMSProvider, func(logger *logger.Logger, conf cryptoengines.CryptoEngineConfig) (cryptoengines.CryptoEngine, error) {
		ceConfig, _ := config.DecodeStruct[AWSCryptoEngine](conf.Config)

		awsCfg, err := aws.GetAwsSdkConfig(ceConfig.AWSSDKConfig)
		if err != nil {
			log.Warnf("skipping AWS KMS engine with id %s: %s", conf.ID, err)
		}

		return NewAWSSecretManagerEngine(logger, *awsCfg, conf.Metadata)
	})
}
