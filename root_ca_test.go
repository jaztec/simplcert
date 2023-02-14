package manager

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"strings"
	"testing"
)

func TestCreateRootCertificate(t *testing.T) {
	t.Run("Should create a root certificate", func(t *testing.T) {
		tests := []struct {
			name     string
			certType CertType
		}{
			{"Test ECDSA root certificate", TypeECDSA},
			{"Test RSA root certificate", TypeRSA},
			{"Test ED25519 root certificate", TypeED25519},
		}

		for _, test := range tests {
			crt, _, err := createRootCertificate(test.certType)
			if err != nil {
				t.Fatalf("%s: Error creating certificate: %+v", test.name, err)
			}

			if _, err := crt.Verify(x509.VerifyOptions{
				KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			}); err != nil {
				// we expect this error, only notify when different error is present
				if strings.Index(err.Error(), "x509: certificate signed by unknown authority") != 0 {
					t.Errorf("%s: Error verifying certificate: %+v", test.name, err)
				}
			}
		}
	})

	t.Run("Should create a leaf certificate", func(t *testing.T) {
		tests := []struct {
			name     string
			certType CertType
		}{
			{"Test ECDSA", TypeECDSA},
			{"Test RSA", TypeRSA},
			{"Test ED25519", TypeED25519},
		}

		for _, test := range tests {
			crt, priv, err := createRootCertificate(test.certType)
			if err != nil {
				t.Fatalf("%s: Error when creating root certificate: %+v", test.name, err)
			}
			p := x509.NewCertPool()
			p.AddCert(crt)
			m := Manager{
				certsPath: "",
				caPriv:    priv,
				caRoot:    crt,
				caPool:    p,
			}
			pem, key, pub, err := m.CreateNamedCert(CertConfig{
				Name:     "Test",
				Host:     "test.org",
				IsServer: false,
			})
			if err != nil {
				t.Fatalf("%s: Error generating cert: %+v", test.name, err)
			}

			priv, err = loadPrivateKey(key)
			if err != nil {
				t.Fatalf("%s: Error loading private key: %+v", test.name, err)
			}

			publ, err := loadPublicKey(pub)
			if err != nil {
				t.Fatalf("%s: Error loading public key: %+v", test.name, err)
			}

			if !publicKeyEquals(priv.Public(), publ) {
				t.Errorf("%s: Error validating public key equals that of private", test.name)
			}

			cert, err := loadCert(pem)
			if err != nil {
				t.Fatalf("%s: Error loading cert: %+v", test.name, err)
			}
			if _, err := cert.Verify(x509.VerifyOptions{
				DNSName:   "test.org",
				KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
				Roots:     m.CaPool(),
			}); err != nil {
				t.Errorf("%s: Error validating cert: %+v", test.name, err)
			}
		}

	})
}

// Promise from `crypto` library:
//
// PublicKey represents a public key using an unspecified algorithm.
//
// Although this type is an empty interface for backwards compatibility reasons,
// all public key types in the standard library implement the following interface
//
//	interface{
//	    Equal(x crypto.PublicKey) bool
//	}
//
// which can be used for increased type safety within applications.
func publicKeyEquals(public crypto.PublicKey, other crypto.PublicKey) bool {
	switch public.(type) {
	case *rsa.PublicKey:
		p := public.(*rsa.PublicKey)
		return p.Equal(other)
	case *ecdsa.PublicKey:
		p := public.(*ecdsa.PublicKey)
		return p.Equal(other)
	case *ed25519.PublicKey:
		p := public.(*ed25519.PublicKey)
		return p.Equal(other)
	}
	return false
}

type EqualKey interface {
	Equal(x crypto.PublicKey) bool
}
