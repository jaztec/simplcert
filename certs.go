package simplcert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"
)

type CertType int
type pemEncoding = string

const (
	rootCABaseFilename             = "root-ca"
	defaultKeySize                 = 2048
	certificatePEM     pemEncoding = "CERTIFICATE"
	privateKeyPEM      pemEncoding = "PRIVATE KEY"
	publicKeyPEM       pemEncoding = "PUBLIC KEY"
)

const (
	TypeECDSA CertType = iota
	TypeRSA
	TypeED25519
)

type CertConfig struct {
	Name         string
	Host         string
	usage        x509.KeyUsage
	extUsage     []x509.ExtKeyUsage
	IsCA         bool
	IsServer     bool
	Country      string
	Organization string
	OutputPath   string
	OutputName   string
	CertType     CertType
	NotAfter     time.Time
}

func createCert(template, parent *x509.Certificate, pub crypto.PublicKey, priv crypto.Signer) (*x509.Certificate, error) {
	crtBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert: %w", err)
	}

	crt, err := x509.ParseCertificate(crtBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cert: %w", err)
	}

	return crt, nil
}

func createNamedCert(cfg CertConfig, parent *x509.Certificate, pub crypto.PublicKey, priv crypto.Signer, objectIdentifier asn1.ObjectIdentifier) (*x509.Certificate, error) {
	san := pkix.Extension{}
	san.Id = objectIdentifier
	san.Critical = false
	san.Value = []byte(fmt.Sprintf("CN=%s", cfg.Name))

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{cfg.Country},
			Organization: []string{cfg.Organization},
		},

		Extensions: []pkix.Extension{
			san,
		},

		NotBefore: time.Now().Add(-10 * time.Second),
		NotAfter:  cfg.NotAfter,

		KeyUsage:    cfg.usage,
		ExtKeyUsage: cfg.extUsage,

		BasicConstraintsValid: true,
		IsCA:                  cfg.IsCA,
	}
	if cfg.Host != "" {
		names := strings.Split(cfg.Host, ",")
		template.DNSNames = make([]string, 0)
		template.IPAddresses = make([]net.IP, 0)
		for _, name := range names {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}
			// check if the hostname is an IP address
			addr := net.ParseIP(name)
			if addr == nil {
				template.DNSNames = append(template.DNSNames, name)
			} else {
				template.IPAddresses = append(template.IPAddresses, addr)
			}
		}
	}
	if parent == nil {
		parent = &template
	}
	crt, err := createCert(&template, parent, pub, priv)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key %w", err)
	}
	return crt, nil
}

func loadPublicKey(pemBytes []byte) (crypto.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(pemBytes))
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch key.(type) {
	case *ecdsa.PublicKey:
		return key.(*ecdsa.PublicKey), err
	case *rsa.PublicKey:
		return key.(*rsa.PublicKey), err
	case ed25519.PublicKey:
		return key.(ed25519.PublicKey), err
	}
	return nil, fmt.Errorf("cert type %T is not a valid type", key)
}

func loadPrivateKey(pemBytes []byte) (crypto.Signer, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(pemBytes))
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	switch key.(type) {
	case *ecdsa.PrivateKey:
		return key.(*ecdsa.PrivateKey), err
	case *rsa.PrivateKey:
		return key.(*rsa.PrivateKey), err
	case ed25519.PrivateKey:
		return key.(ed25519.PrivateKey), err
	}
	return nil, fmt.Errorf("cert type %T is not a valid type", key)
}

func loadCert(bytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, fmt.Errorf("decoding file from %s failed", string(bytes))
	}

	return x509.ParseCertificate(block.Bytes)
}

func createKey(t CertType) (crypto.Signer, error) {
	switch t {
	case TypeECDSA:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case TypeRSA:
		return rsa.GenerateKey(rand.Reader, defaultKeySize)
	case TypeED25519:
		_, priv, err := ed25519.GenerateKey(rand.Reader)
		return priv, err
	default:
		return nil, fmt.Errorf("type %d is not a valid CertType", t)
	}
}

func EncodePublicKey(rawBytes []byte) []byte {
	return encodePEM(publicKeyPEM, rawBytes)
}

func EncodePrivateKey(rawBytes []byte) []byte {
	return encodePEM(privateKeyPEM, rawBytes)
}

func EncodeCertificate(rawBytes []byte) []byte {
	return encodePEM(certificatePEM, rawBytes)
}

func encodePEM(t pemEncoding, rawBytes []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:  t,
		Bytes: rawBytes,
	})
}
