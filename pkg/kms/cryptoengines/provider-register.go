package cryptoengines

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"

var cryptoEngineBuilders = make(map[CryptoEngineProvider]func(*logger.Logger, CryptoEngineConfig) (CryptoEngine, error))

func RegisterCryptoEngine(name CryptoEngineProvider, builder func(*logger.Logger, CryptoEngineConfig) (CryptoEngine, error)) {
	cryptoEngineBuilders[name] = builder
}

func GetEngineBuilder(name CryptoEngineProvider) func(*logger.Logger, CryptoEngineConfig) (CryptoEngine, error) {
	return cryptoEngineBuilders[name]
}
