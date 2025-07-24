package cryptoengines

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"

type CryptoEngineProvider string

const (
	HashicorpVaultProvider    CryptoEngineProvider = "hashicorp_vault"
	AWSKMSProvider            CryptoEngineProvider = "aws_kms"
	AWSSecretsManagerProvider CryptoEngineProvider = "aws_secrets_manager"
	FilesystemProvider        CryptoEngineProvider = "filesystem"
	PKCS11Provider            CryptoEngineProvider = "pkcs11"
)

type CryptoEngineConfig struct {
	ID       string                 `mapstructure:"id"`
	Metadata map[string]interface{} `mapstructure:"metadata"`
	Type     CryptoEngineProvider   `mapstructure:"type"`
	Config   map[string]interface{} `mapstructure:",remain"`
}

type CryptoEngineConfigAdapter[E any] struct {
	ID       string
	Metadata map[string]interface{}
	Type     CryptoEngineProvider
	Config   E
}

func (c CryptoEngineConfigAdapter[E]) Marshal(ce CryptoEngineConfig) (*CryptoEngineConfigAdapter[E], error) {
	config, err := config.DecodeStruct[E](ce.Config)
	if err != nil {
		return nil, err
	}
	return &CryptoEngineConfigAdapter[E]{
		ID:       ce.ID,
		Metadata: ce.Metadata,
		Type:     ce.Type,
		Config:   config,
	}, nil
}

func (c CryptoEngineConfigAdapter[E]) Unmarshal() (*CryptoEngineConfig, error) {

	config, err := config.EncodeStruct(c.Config)
	if err != nil {
		return nil, err
	}

	return &CryptoEngineConfig{
		ID:       c.ID,
		Metadata: c.Metadata,
		Type:     c.Type,
		Config:   config,
	}, nil
}
