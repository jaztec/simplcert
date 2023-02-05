package security

import (
	"crypto/x509"
	"strings"
	"testing"
)

func TestCreateRootCertificate(t *testing.T) {
	t.Run("Should create a root certificate", func(t *testing.T) {
		crt, p, k, err := createRootCertificate()
		if err != nil {
			t.Fatalf("Error creating certificate: %+v", err)
		}

		expectSize := 1127
		if len(p) != expectSize {
			t.Errorf("Expected PEM to be of size %d but got %d", expectSize, len(p))
		}

		if err := k.Validate(); err != nil {
			t.Errorf("Private key not valid: %+v", err)
		}

		if _, err := crt.Verify(x509.VerifyOptions{
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		}); err != nil {
			// we expect this error, only notify when different error is present
			if strings.Index(err.Error(), "x509: certificate signed by unknown authority") != 0 {
				t.Errorf("Error verifying certificate: %+v", err)
			}
		}
	})

	t.Run("Should create a leaf certificate", func(t *testing.T) {
		crt, _, priv, _ := createRootCertificate()
		p := x509.NewCertPool()
		p.AddCert(crt)
		m := Manager{
			certsPath: "",
			caPriv:    priv,
			caRoot:    crt,
			caPool:    p,
		}
		pem, key, pub, err := m.CreateNamedCert("Test", "test.org", false)
		if err != nil {
			t.Fatalf("Error generating cert: %+v", err)
		}

		priv, err = loadPrivateKey(key)
		if err != nil {
			t.Fatalf("Error loading private key: %+v", err)
		}
		if err := priv.Validate(); err != nil {
			t.Errorf("Error validating private key: %+v", err)
		}
		publ, err := loadPublicKey(pub)
		if err != nil {
			t.Fatalf("Error loading public key: %+v", err)
		}

		if !priv.PublicKey.Equal(publ) {
			t.Error("Error validating public key equals that of private")
		}

		cert, err := loadCert(pem)
		if err != nil {
			t.Fatalf("Error loading cert: %+v", err)
		}
		if _, err := cert.Verify(x509.VerifyOptions{
			DNSName:   "test.org",
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			Roots:     m.CaPool(),
		}); err != nil {
			t.Errorf("Error validating cert: %+v", err)
		}

	})
}
