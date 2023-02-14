package simplcert

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	log "github.com/sirupsen/logrus"
	math_rand "math/rand"
	"time"
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

func (m *Manager) RootCrt() *x509.Certificate {
	return m.caRoot
}

// CreateNamedCert will return raw TLS certificate, Private key and Public key bytes
func (m *Manager) CreateNamedCert(cfg CertConfig) (*x509.Certificate, crypto.Signer, []byte, error) {
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

	math_rand.Seed(time.Now().UnixMilli())
	identifier := baseIdentifier
	for i := 0; i < 3; i++ {
		identifier = append(identifier, math_rand.Intn(32))
	}

	crt, err := createNamedCert(cfg, m.caRoot, &priv.PublicKey, m.caPriv, identifier)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}

	p, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	pub := EncodePublicKey(p)

	return crt, priv, pub, nil
}

func (m *Manager) MarshalPrivateKey(key crypto.PrivateKey) ([]byte, error) {
	b, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate security keys: %w", err)
	}
	return b, nil
}

func NewManager(certsPath string) (*Manager, error) {
	crt, priv, err := loadCA(certsPath)
	if err != nil {
		log.WithField("error", err).Debug("Error loading CA files")
		return nil, fmt.Errorf("no root certificates found at %s, please generate them first: %w", certsPath, NoCertsError)
	}

	p := x509.NewCertPool()
	p.AddCert(crt)

	m := &Manager{certsPath: certsPath, caPriv: priv, caRoot: crt, caPool: p}

	return m, nil
}
