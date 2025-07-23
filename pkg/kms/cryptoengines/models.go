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

type CryptoEngineType string
type SupportedKeyTypeInfo struct {
	Type  string `json:"type"`
	Sizes []int  `json:"sizes"`
}
