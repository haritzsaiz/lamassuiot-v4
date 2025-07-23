package cryptoengines

type CryptoEngineProvider string

type CryptoEngineConfig struct {
	ID       string                 `mapstructure:"id"`
	Metadata map[string]interface{} `mapstructure:"metadata"`
	Type     CryptoEngineProvider   `mapstructure:"type"`
	Config   map[string]interface{} `mapstructure:",remain"`
}
