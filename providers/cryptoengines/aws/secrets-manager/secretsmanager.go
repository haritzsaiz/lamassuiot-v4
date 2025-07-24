package secretsmanager

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/lamassuiot/lamassuiot/v4/pkg/kms/cryptoengines"
	httpclient "github.com/lamassuiot/lamassuiot/v4/pkg/shared/http/client"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

type AWSSecretsManagerCryptoEngine struct {
	config      cryptoengines.CryptoEngineInfo
	smngerCli   *secretsmanager.Client
	logger      *logger.Logger
	keyProvider *cryptoengines.SoftwareKeyProvider
}

func NewAWSSecretManagerEngine(logger *logger.Logger, awsConf aws.Config, metadata map[string]any) (cryptoengines.CryptoEngine, error) {
	lAWSSM := logger.With("subsystem-provider", "AWS SecretsManager Client")

	httpCli, err := httpclient.BuildHTTPClientWithTracerLogger(http.DefaultClient, lAWSSM)
	if err != nil {
		return nil, err
	}

	awsConf.HTTPClient = httpCli

	smCli := secretsmanager.NewFromConfig(awsConf)

	return &AWSSecretsManagerCryptoEngine{
		logger:      lAWSSM,
		smngerCli:   smCli,
		keyProvider: cryptoengines.NewSoftwareKeyProvider(logger),
		config: cryptoengines.CryptoEngineInfo{
			Type:          "AWSSecretsManager",
			SecurityLevel: cryptoengines.SL1,
			Provider:      "Amazon Web Services",
			Name:          "Secrets Manager",
			Metadata:      metadata,
			SupportedKeyTypes: []cryptoengines.SupportedKeyTypeInfo{
				{
					Type: "RSA",
					Sizes: []int{
						2048,
						3072,
						4096,
					},
				},
				{
					Type: "ECDSA",
					Sizes: []int{
						224,
						256,
						521,
					},
				},
			},
		},
	}, nil
}

func (engine *AWSSecretsManagerCryptoEngine) GetEngineConfig(ctx context.Context) cryptoengines.CryptoEngineInfo {
	return engine.config
}

func (engine *AWSSecretsManagerCryptoEngine) GetPrivateKeyByID(ctx context.Context, keyID string) (crypto.Signer, error) {
	engine.logger.Debugf("Getting the private key with ID: %s", keyID)

	result, err := engine.smngerCli.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(keyID),
	})
	if err != nil {
		engine.logger.Errorf("could not get Secret Value: %s", err)
		return nil, err
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	var keyMap map[string]string

	err = json.Unmarshal([]byte(secretString), &keyMap)
	if err != nil {
		return nil, err
	}

	pemBytes, ok := keyMap["key"]
	if !ok {
		engine.logger.Errorf("'key' variable not found in secret")
		return nil, fmt.Errorf("'key' not found in secret")
	}

	decodedPemBytes, err := base64.StdEncoding.DecodeString(pemBytes)
	if err != nil {
		engine.logger.Errorf("could not decode key: %s", err)
		return nil, err
	}

	return engine.keyProvider.ParsePrivateKey(decodedPemBytes)
}

func (engine *AWSSecretsManagerCryptoEngine) ListPrivateKeyIDs(ctx context.Context) ([]string, error) {
	engine.logger.Debugf("listing private key IDs")

	keyRes, err := engine.smngerCli.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
	if err != nil {
		engine.logger.Errorf("could not list secrets: %s", err)
		return nil, err
	}

	keys := []string{}
	for _, secret := range keyRes.SecretList {
		keys = append(keys, *secret.Name)
	}

	engine.logger.Debugf("private key IDs successfully listed")

	return keys, nil
}

func (engine *AWSSecretsManagerCryptoEngine) CreateRSAPrivateKey(ctx context.Context, keySize int) (string, crypto.Signer, error) {
	engine.logger.Debugf("creating RSA private key")

	_, key, err := engine.keyProvider.CreateRSAPrivateKey(keySize)
	if err != nil {
		engine.logger.Errorf("could not create RSA private key: %s", err)
		return "", nil, err
	}

	engine.logger.Debugf("RSA key successfully generated")
	return engine.importKey(ctx, key)
}

func (engine *AWSSecretsManagerCryptoEngine) CreateECDSAPrivateKey(ctx context.Context, curve elliptic.Curve) (string, crypto.Signer, error) {
	engine.logger.Debugf("creating ECDSA private key")

	_, key, err := engine.keyProvider.CreateECDSAPrivateKey(curve)
	if err != nil {
		engine.logger.Errorf("could not create ECDSA private key: %s", err)
		return "", nil, err
	}

	engine.logger.Debugf("ECDSA key successfully generated")
	return engine.importKey(ctx, key)
}

func (engine *AWSSecretsManagerCryptoEngine) ImportRSAPrivateKey(ctx context.Context, key *rsa.PrivateKey) (string, crypto.Signer, error) {
	engine.logger.Debugf("importing RSA private key")

	keyID, signer, err := engine.importKey(ctx, key)
	if err != nil {
		engine.logger.Errorf("could not import RSA key: %s", err)
		return "", nil, err
	}

	engine.logger.Debugf("RSA key successfully imported")
	return keyID, signer, nil
}

func (engine *AWSSecretsManagerCryptoEngine) ImportECDSAPrivateKey(ctx context.Context, key *ecdsa.PrivateKey) (string, crypto.Signer, error) {
	engine.logger.Debugf("importing ECDSA private key")

	keyID, signer, err := engine.importKey(ctx, key)
	if err != nil {
		engine.logger.Errorf("could not import ECDSA key: %s", err)
		return "", nil, err
	}

	engine.logger.Debugf("ECDSA key successfully imported")
	return keyID, signer, nil
}

func (engine *AWSSecretsManagerCryptoEngine) importKey(ctx context.Context, key crypto.Signer) (string, crypto.Signer, error) {
	pubKey := key.Public()

	keyID, err := engine.keyProvider.EncodePKIXPublicKeyDigest(pubKey)
	if err != nil {
		engine.logger.Errorf("could not encode public key digest: %s", err)
		return "", nil, err
	}

	b64PemKey, err := engine.keyProvider.MarshalAndEncodePKIXPrivateKey(key)
	if err != nil {
		engine.logger.Errorf("could not marshal and encode private key: %s", err)
		return "", nil, err
	}

	keyVal := `{"key": "` + b64PemKey + `"}`

	_, err = engine.smngerCli.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(keyID),
		SecretString: aws.String(keyVal),
	})

	if err != nil {
		engine.logger.Error("Could not import private key: ", err)
		return "", nil, err
	}

	return keyID, key, nil
}

func (engine *AWSSecretsManagerCryptoEngine) RenameKey(ctx context.Context, oldID, newID string) error {
	engine.logger.Debugf("renaming key with ID: %s to %s", oldID, newID)

	result, err := engine.smngerCli.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(oldID),
	})
	if err != nil {
		engine.logger.Errorf("could not get Secret Value: %s", err)
		return err
	}

	_, err = engine.smngerCli.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Name:         aws.String(newID),
		SecretString: result.SecretString,
	})
	if err != nil {
		engine.logger.Errorf("could not create Secret Value: %s", err)
		return err
	}

	err = engine.DeleteKey(ctx, oldID)
	if err != nil {
		engine.logger.Errorf("could not delete old key: %s", err)
	}

	engine.logger.Debugf("key successfully renamed")
	return nil
}

func (engine *AWSSecretsManagerCryptoEngine) DeleteKey(ctx context.Context, keyID string) error {
	engine.logger.Debugf("deleting key with ID: %s", keyID)

	_, err := engine.smngerCli.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		SecretId:             aws.String(keyID),
		RecoveryWindowInDays: aws.Int64(7),
	})

	if err != nil {
		engine.logger.Errorf("could not delete key: %s", err)
		return err
	}

	engine.logger.Debugf("key successfully deleted")
	return nil
}
