package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

const (
	rootCABaseFilename = "root-ca"
)

type CertConfig struct {
	Name     string
	Host     string
	Usage    x509.KeyUsage
	ExtUsage []x509.ExtKeyUsage
	IsCA     bool
}

func createNamedCert(cfg CertConfig, parent *x509.Certificate, pub *rsa.PublicKey, priv *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	san := pkix.Extension{}
	san.Id = asn1.ObjectIdentifier{2, 5, 29, 17}
	san.Critical = false
	san.Value = []byte(fmt.Sprintf("CN=%s", cfg.Name))

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{"NL"},
			Organization: []string{"SERP"},
			//CommonName:   cfg.Name,
		},

		Extensions: []pkix.Extension{
			san,
		},

		NotBefore: time.Now().Add(-10 * time.Second),
		NotAfter:  time.Now().AddDate(10, 0, 0),

		KeyUsage:    cfg.Usage,
		ExtKeyUsage: cfg.ExtUsage,

		BasicConstraintsValid: true,
		IsCA:                  cfg.IsCA,
	}
	if cfg.Host != "" {
		template.DNSNames = []string{cfg.Host}
	}
	if parent == nil {
		parent = &template
	}
	crt, crtPem, err := createCert(&template, parent, pub, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key %w", err)
	}
	return crt, crtPem, nil
}

func createCert(template, parent *x509.Certificate, pub *rsa.PublicKey, priv *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	crtBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cert: %w", err)
	}

	crt, err := x509.ParseCertificate(crtBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create cert: %w", err)
	}

	crtPem := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crtBytes,
	})

	return crt, crtPem, nil
}

func loadPublicKey(bytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(bytes))
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func loadPrivateKey(bytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(bytes))
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	return key.(*rsa.PrivateKey), err
}

func loadCert(bytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(bytes))
	}

	return x509.ParseCertificate(block.Bytes)
}
