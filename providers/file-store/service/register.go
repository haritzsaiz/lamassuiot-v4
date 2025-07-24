package service

import (
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"gocloud.dev/blob"
)

var fileStoreBuilders = make(map[FileStoreProvider]func(*logger.Logger, FileStoreConfig) (*blob.Bucket, error))

func RegisterProvider(name FileStoreProvider, builder func(*logger.Logger, FileStoreConfig) (*blob.Bucket, error)) {
	fileStoreBuilders[name] = builder
}

func GetProvider(name FileStoreProvider) func(*logger.Logger, FileStoreConfig) (*blob.Bucket, error) {
	return fileStoreBuilders[name]
}
