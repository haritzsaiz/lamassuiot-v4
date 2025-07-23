package cryptoengines

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
)

type SoftwareKeyProvider struct {
	logger *logger.Logger
}

func NewSoftwareKeyProvider(logger *logger.Logger) *SoftwareKeyProvider {
	return &SoftwareKeyProvider{
		logger: logger,
	}
}

// CreateRSAPrivateKey creates a RSA private key with the specified key size
func (p *SoftwareKeyProvider) CreateRSAPrivateKey(keySize int) (string, *rsa.PrivateKey, error) {
	lFunc := p.logger

	lFunc.Infof("starting RSA key generation with %d bit key size", keySize)

	// Validate key size
	if keySize < 2048 {
		lFunc.Warnf("RSA key size %d is below recommended minimum of 2048 bits", keySize)
	}

	lFunc.Debugf("generating RSA private key using crypto/rand reader")
	key, err := rsa.GenerateKey(rand.Reader, keySize)

	if err != nil {
		lFunc.Errorf("RSA key generation failed for %d bit key: %s", keySize, err)
		return "", nil, err
	}

	lFunc.Debugf("RSA key generation successful - key size: %d bits, public exponent: %d",
		key.Size()*8, key.PublicKey.E)

	lFunc.Debugf("encoding public key digest for RSA key")
	encDigest, err := p.EncodePKIXPublicKeyDigest(&key.PublicKey)
	if err != nil {
		lFunc.Errorf("failed to encode public key digest for RSA key: %s", err)
		return "", nil, err
	}

	lFunc.Infof("RSA key creation completed successfully - digest: %s", encDigest)
	return encDigest, key, nil
}

func (p *SoftwareKeyProvider) CreateECDSAPrivateKey(curve elliptic.Curve) (string, *ecdsa.PrivateKey, error) {
	lFunc := p.logger

	curveName := curve.Params().Name
	lFunc.Infof("starting ECDSA key generation with curve: %s", curveName)
	lFunc.Debugf("curve parameters - name: %s, bit size: %d", curveName, curve.Params().BitSize)

	lFunc.Debugf("generating ECDSA private key using crypto/rand reader")
	key, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		lFunc.Errorf("ECDSA key generation failed for curve %s: %s", curveName, err)
		return "", nil, err
	}

	lFunc.Debugf("ECDSA key generation successful - curve: %s, private key size: %d bits",
		curveName, key.Curve.Params().BitSize)

	lFunc.Debugf("encoding public key digest for ECDSA key")
	encDigest, err := p.EncodePKIXPublicKeyDigest(&key.PublicKey)
	if err != nil {
		lFunc.Errorf("failed to encode public key digest for ECDSA key: %s", err)
		return "", nil, err
	}

	lFunc.Infof("ECDSA key creation completed successfully - curve: %s, digest: %s", curveName, encDigest)
	return encDigest, key, nil
}

func (p *SoftwareKeyProvider) MarshalAndEncodePKIXPrivateKey(key interface{}) (string, error) {
	p.logger.Infof("starting private key marshaling and encoding process")

	// Log key type information
	switch k := key.(type) {
	case *rsa.PrivateKey:
		p.logger.Debugf("marshaling RSA private key - key size: %d bits", k.Size()*8)
	case *ecdsa.PrivateKey:
		p.logger.Debugf("marshaling ECDSA private key - curve: %s", k.Curve.Params().Name)
	default:
		p.logger.Debugf("marshaling private key of type: %T", key)
	}

	p.logger.Debugf("marshaling private key to PKCS#8 format")
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		p.logger.Errorf("PKCS#8 marshaling failed: %s", err)
		return "", err
	}

	p.logger.Debugf("PKCS#8 marshaling successful - key data size: %d bytes", len(keyBytes))

	p.logger.Debugf("encoding private key to PEM format")
	keyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	p.logger.Debugf("PEM encoding successful - PEM data size: %d bytes", len(keyPem))

	keyBase64 := base64.StdEncoding.EncodeToString([]byte(keyPem))
	p.logger.Infof("private key marshaling and encoding completed successfully")
	return keyBase64, nil
}

func (p *SoftwareKeyProvider) EncodePKIXPublicKeyDigest(key interface{}) (string, error) {
	p.logger.Debugf("starting public key digest extraction and encoding")

	// Log key type information
	switch k := key.(type) {
	case *rsa.PublicKey:
		p.logger.Debugf("processing RSA public key - key size: %d bits, exponent: %d", k.Size()*8, k.E)
	case *ecdsa.PublicKey:
		p.logger.Debugf("processing ECDSA public key - curve: %s", k.Curve.Params().Name)
	default:
		p.logger.Debugf("processing public key of type: %T", key)
	}

	p.logger.Debugf("marshaling public key to PKIX format")
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		p.logger.Errorf("PKIX public key marshaling failed: %s", err)
		return "", err
	}

	p.logger.Debugf("PKIX marshaling successful - public key data size: %d bytes", len(pubkeyBytes))

	p.logger.Debugf("computing SHA-256 hash of public key")
	hash := sha256.New()
	hash.Write(pubkeyBytes)
	digest := hash.Sum(nil)
	p.logger.Tracef("SHA-256 digest (raw bytes): %x", digest)

	p.logger.Debugf("encoding digest to hexadecimal string")
	hexDigest := hex.EncodeToString(digest)
	p.logger.Infof("public key digest computed successfully: %s", hexDigest)

	return hexDigest, nil
}

func (p *SoftwareKeyProvider) ParsePrivateKey(pemBytes []byte) (crypto.Signer, error) {
	p.logger.Infof("starting private key parsing from PEM data")
	p.logger.Debugf("input PEM data size: %d bytes", len(pemBytes))

	p.logger.Debugf("decoding PEM block")
	block, remainder := pem.Decode(pemBytes)
	if block == nil {
		p.logger.Errorf("PEM decoding failed - no valid PEM block found")
		return nil, fmt.Errorf("no key found")
	}

	p.logger.Debugf("PEM block decoded successfully - type: %s, data size: %d bytes",
		block.Type, len(block.Bytes))

	if len(remainder) > 0 {
		p.logger.Debugf("additional data found after PEM block: %d bytes", len(remainder))
	}

	var genericKey interface{}
	var err error
	var keyFormat string

	// First try to parse as PKCS8
	p.logger.Debugf("attempting to parse as PKCS#8 private key")
	genericKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		p.logger.Debugf("PKCS#8 parsing failed: %s", err)

		// If it fails, try to parse as PKCS1
		p.logger.Debugf("attempting to parse as PKCS#1 RSA private key")
		genericKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			p.logger.Debugf("PKCS#1 parsing failed: %s", err)

			// If it fails, try to parse as EC
			p.logger.Debugf("attempting to parse as EC private key")
			genericKey, err = x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				p.logger.Errorf("all private key parsing attempts failed - PKCS#8, PKCS#1, and EC: %s", err)
				return nil, err
			}
			keyFormat = "EC"
		} else {
			keyFormat = "PKCS#1 RSA"
		}
	} else {
		keyFormat = "PKCS#8"
	}

	p.logger.Infof("private key parsed successfully using %s format", keyFormat)

	switch key := genericKey.(type) {
	case *rsa.PrivateKey:
		p.logger.Infof("parsed RSA private key - key size: %d bits, public exponent: %d",
			key.Size()*8, key.PublicKey.E)
		return key, nil
	case *ecdsa.PrivateKey:
		p.logger.Infof("parsed ECDSA private key - curve: %s, bit size: %d",
			key.Curve.Params().Name, key.Curve.Params().BitSize)
		return key, nil
	default:
		p.logger.Errorf("unsupported private key type: %T", key)
		return nil, errors.New("unsupported key type")
	}
}
