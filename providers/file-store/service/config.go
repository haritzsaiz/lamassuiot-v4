package service

import "github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"

type FileStoreProvider string

const (
	LocalFilesystem FileStoreProvider = "local"
	AWSS3           FileStoreProvider = "s3"
)

type FileStoreConfig struct {
	ID       string                 `mapstructure:"id"`
	Metadata map[string]interface{} `mapstructure:"metadata"`
	Type     FileStoreProvider      `mapstructure:"type"`
	Config   map[string]interface{} `mapstructure:",remain"`
}

type FileStoreConfigAdapter[E any] struct {
	ID       string
	Metadata map[string]interface{}
	Type     FileStoreProvider
	Config   E
}

func (c FileStoreConfigAdapter[E]) Marshal(ce FileStoreConfig) (*FileStoreConfigAdapter[E], error) {
	config, err := config.DecodeStruct[E](ce.Config)
	if err != nil {
		return nil, err
	}
	return &FileStoreConfigAdapter[E]{
		ID:       ce.ID,
		Metadata: ce.Metadata,
		Type:     ce.Type,
		Config:   config,
	}, nil
}

func (c FileStoreConfigAdapter[E]) Unmarshal() (*FileStoreConfig, error) {

	config, err := config.EncodeStruct(c.Config)
	if err != nil {
		return nil, err
	}

	return &FileStoreConfig{
		ID:       c.ID,
		Metadata: c.Metadata,
		Type:     c.Type,
		Config:   config,
	}, nil
}
