package cryptoengines

type CryptoEngineInfo struct {
	Type              CryptoEngineType       `json:"type"`
	SecurityLevel     CryptoEngineSL         `json:"security_level"`
	Provider          string                 `json:"provider"`
	Name              string                 `json:"name"`
	Metadata          map[string]any         `json:"metadata"`
	SupportedKeyTypes []SupportedKeyTypeInfo `json:"supported_key_types"`
}

type CryptoEngineSL int

const (
	SL0 CryptoEngineSL = 0
	SL1 CryptoEngineSL = 1
	SL2 CryptoEngineSL = 2
)

type SupportedKeyTypeInfo struct {
	Type  string `json:"type"`
	Sizes []int  `json:"sizes"`
}

type CryptoEngineType string

const (
	PKCS11            CryptoEngineType = "PKCS11"
	AzureKeyVault     CryptoEngineType = "AZURE_KEY_VAULT"
	VaultKV2          CryptoEngineType = "HASHICORP_VAULT_KV_V2"
	AWSKMS            CryptoEngineType = "AWS_KMS"
	AWSSecretsManager CryptoEngineType = "AWS_SECRETS_MANAGER"
	Filesystem        CryptoEngineType = "FILESYSTEM"
)
