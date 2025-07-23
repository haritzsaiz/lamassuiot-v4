package cryptoengines

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
)

type CryptoEngine interface {
	GetEngineConfig(context.Context) CryptoEngineInfo

	ListPrivateKeyIDs(context.Context) ([]string, error)
	GetPrivateKeyByID(context.Context, string) (crypto.Signer, error)

	CreateRSAPrivateKey(context.Context, int) (string, crypto.Signer, error)
	CreateECDSAPrivateKey(context.Context, elliptic.Curve) (string, crypto.Signer, error)

	ImportRSAPrivateKey(ctx context.Context, key *rsa.PrivateKey) (string, crypto.Signer, error)
	ImportECDSAPrivateKey(ctx context.Context, key *ecdsa.PrivateKey) (string, crypto.Signer, error)

	DeleteKey(ctx context.Context, keyID string) error

	RenameKey(ctx context.Context, oldID, newID string) error
}
