package manager

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
)

type internalError string

func (i internalError) Error() string {
	return string(i)
}

const NoCertsError = internalError("no root certificates found")

type Manager struct {
	certsPath string
	caPriv    crypto.Signer
	caRoot    *x509.Certificate
	caPool    *x509.CertPool
}

func (m *Manager) CaPool() *x509.CertPool {
	return m.caPool
}

func (m *Manager) RootPEM() []byte {
	return EncodeCertificate(m.caRoot.Raw)
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

	crt, err := createNamedCert(cfg, m.caRoot, &priv.PublicKey, m.caPriv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}

	b, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	key = EncodePrivateKey(b)

	p, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	pub = EncodePublicKey(p)

	crtPem = EncodeCertificate(crt.Raw)

	return crtPem, key, pub, nil
}

func NewManager(certsPath string) (*Manager, error) {
	crt, priv, err := loadCA(certsPath)
	if err != nil {
		return nil, fmt.Errorf("no root certificates found at %s, please generate them first: %w", certsPath, NoCertsError)
	}

	p := x509.NewCertPool()
	p.AddCert(crt)

	m := &Manager{certsPath: certsPath, caPriv: priv, caRoot: crt, caPool: p}

	return m, nil
}
