package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type Manager struct {
	certsPath string
	caPriv    *ecdsa.PrivateKey
	caRoot    *x509.Certificate
	caPool    *x509.CertPool
	rootPem   []byte
}

func (m *Manager) CaPool() *x509.CertPool {
	return m.caPool
}

func (m *Manager) RootPEM() []byte {
	return m.rootPem
}

// CreateNamedCert will return raw TLS certificate, Private key and Public key bytes
func (m *Manager) CreateNamedCert(cfg CertConfig) (crtPem []byte, key []byte, pub []byte, err error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate key")
	}
	if cfg.IsServer {
		cfg.usage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		cfg.extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		cfg.usage = x509.KeyUsageContentCommitment
		cfg.extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}

	crt, _, err := createNamedCert(cfg, m.caRoot, &priv.PublicKey, m.caPriv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}

	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	key = pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: b,
		},
	)

	p, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	pub = pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: p,
		},
	)

	crtPem = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: crt.Raw,
	})

	return crtPem, key, pub, nil
}

func NewManager(certsPath string) (*Manager, error) {
	crt, priv, crtPem, err := loadCA(certsPath)
	if err != nil {
		if err = createRootCAFiles(certsPath); err != nil {
			return nil, err
		}
		if crt, priv, crtPem, err = loadCA(certsPath); err != nil {
			return nil, err
		}
	}

	p := x509.NewCertPool()
	p.AddCert(crt)

	m := &Manager{certsPath: certsPath, caPriv: priv, caRoot: crt, caPool: p, rootPem: crtPem}

	return m, nil
}
