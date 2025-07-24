package cryptoutils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

// ReadCertificateFromFile reads and parses an X.509 certificate from a file
func ReadCertificateFromFile(filePath string) (*x509.Certificate, error) {
	if filePath == "" {
		return nil, fmt.Errorf("cannot open empty filepath")
	}

	certFileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParseCertificate(string(certFileBytes))
}

// ParseCertificate parses an X.509 certificate from PEM-encoded string
func ParseCertificate(cert string) (*x509.Certificate, error) {
	certDERBlock, _ := pem.Decode([]byte(cert))
	if certDERBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM certificate")
	}
	return x509.ParseCertificate(certDERBlock.Bytes)
}

// ParseCertificateRequest parses an X.509 certificate signing request from PEM-encoded string
func ParseCertificateRequest(cert string) (*x509.CertificateRequest, error) {
	certDERBlock, _ := pem.Decode([]byte(cert))
	if certDERBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM certificate request")
	}
	return x509.ParseCertificateRequest(certDERBlock.Bytes)
}

// ReadPrivateKeyFromFile reads and parses a private key from a file
func ReadPrivateKeyFromFile(filePath string) (interface{}, error) {
	keyFileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return ParsePrivateKey(keyFileBytes)
}

// ParsePrivateKey parses a private key from PEM-encoded bytes
// Supports PKCS#1, PKCS#8, and EC private key formats
func ParsePrivateKey(privKeyBytes []byte) (interface{}, error) {
	keyDERBlock, _ := pem.Decode(privKeyBytes)
	if keyDERBlock == nil {
		return nil, fmt.Errorf("failed to decode PEM private key")
	}

	if key, err := x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes); err == nil {
		switch key := key.(type) {
		case *rsa.PrivateKey, *ecdsa.PrivateKey, ed25519.PrivateKey:
			return key, nil
		default:
			return nil, errors.New("tls: found unknown private key type in PKCS#8 wrapping")
		}
	}
	if key, err := x509.ParseECPrivateKey(keyDERBlock.Bytes); err == nil {
		return key, nil
	}

	return nil, errors.New("tls: failed to parse private key")
}

// CertificateToPEM converts an X.509 certificate to PEM-encoded string
func CertificateToPEM(c *x509.Certificate) string {
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.Raw})
	return string(pemCert)
}

// PrivateKeyToPEM converts a private key to PEM-encoded string using PKCS#8 format
func PrivateKeyToPEM(key any) (string, error) {
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return "", err
	}

	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: b,
		},
	)

	return string(pemdata), nil
}

// GenerateSelfSignedCertificate generates a self-signed X.509 certificate for the given key and common name
func GenerateSelfSignedCertificate(key crypto.Signer, cn string) (*x509.Certificate, error) {
	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 160))
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %w", err)
	}

	crt := x509.Certificate{
		SerialNumber: sn,
		Subject:      pkix.Name{CommonName: cn},
	}

	crtB, err := x509.CreateCertificate(rand.Reader, &crt, &crt, key.Public(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	crtP, err := x509.ParseCertificate(crtB)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created certificate: %w", err)
	}

	return crtP, nil
}

func LoadSytemCACertPool() *x509.CertPool {
	certPool := x509.NewCertPool()
	systemCertPool, err := x509.SystemCertPool()
	if err == nil {
		certPool = systemCertPool
	} else {
		log.Warnf("could not get system cert pool (trusted CAs). Using empty pool: %s", err)
	}

	return certPool
}

func LoadSystemCACertPoolWithExtraCAsFromFiles(casToAdd []string) *x509.CertPool {
	certPool := x509.NewCertPool()
	systemCertPool, err := x509.SystemCertPool()
	if err == nil {
		certPool = systemCertPool
	} else {
		log.Warnf("could not get system cert pool (trusted CAs). Using empty pool: %s", err)
	}

	for _, ca := range casToAdd {
		if ca == "" {
			continue
		}

		caCert, err := ReadCertificateFromFile(ca)
		if err != nil {
			log.Warnf("could not load CA certificate in %s. Skipping CA: %s", ca, err)
			continue
		}

		certPool.AddCert(caCert)
	}

	return certPool
}
