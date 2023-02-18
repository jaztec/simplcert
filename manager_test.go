package simplcert_test

import (
	"crypto/x509"
	"errors"
	"github.com/jaztec/simplcert"
	"os"
	"testing"
	"time"
)

func createTestDirectory(name string, t *testing.T) string {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("%s: Cannot get work dir: %+v", name, err)
	}
	outPath := cwd + string(os.PathSeparator) + "tmp"
	if err := os.Mkdir(outPath, 0777); err != nil {
		t.Fatalf("%s: Cannot create work dir: %+v", name, err)
	}
	return outPath
}

func certConfig(certType simplcert.CertType, isServer bool) simplcert.CertConfig {
	return simplcert.CertConfig{
		Name:     "Test",
		Host:     "test.org, , 127.0.0.1, ::1,",
		IsServer: isServer,
		CertType: certType,
		NotAfter: time.Now().AddDate(0, 0, 1),
	}
}

func runGenerateCertTest(m *simplcert.Manager, name string, certType simplcert.CertType, isServer bool, t *testing.T) {
	crt, _, _, err := m.CreateNamedCert(certConfig(certType, isServer))
	if err != nil {
		t.Fatalf("%s (server %t): Error creating named certificate: %+v", name, isServer, err)
	}
	var keyUsage []x509.ExtKeyUsage
	if isServer {
		keyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	} else {
		keyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	}
	if _, err := crt.Verify(x509.VerifyOptions{
		DNSName:   "test.org",
		KeyUsages: keyUsage,
		Roots:     m.CaPool(),
	}); err != nil {
		t.Errorf("%s (server %t): Error validating cert: %+v", name, isServer, err)
	}
}

func TestManagerFunctions(t *testing.T) {
	tests := []struct {
		name     string
		certType simplcert.CertType
	}{
		{"Test manager functions ECDSA", simplcert.TypeECDSA},
		{"Test manager functions RSA", simplcert.TypeRSA},
		{"Test manager functions ED25519", simplcert.TypeED25519},
	}

	for _, test := range tests {
		outPath := createTestDirectory(test.name, t)
		m, err := simplcert.NewManager(outPath)
		if err == nil {
			t.Error("An error should occur when root certs are not yet generated")
		}
		if err != nil && !errors.Is(err, simplcert.NoCertsError) {
			t.Errorf("%s: The return error should be of type %T, got %T", test.name, simplcert.NoCertsError, err)
		}
		if m != nil {
			t.Errorf("%s: The returned manager type should be nil", test.name)
		}
		if err := simplcert.CreateRootCAFiles(test.certType, outPath); err != nil {
			t.Errorf("%s: Creating roto CA files returned an error: %+v", test.name, err)
		}
		m, err = simplcert.NewManager(outPath)
		if err != nil {
			t.Errorf("%s: After generating root CA files the NewManager cannot return an error, got %+v", test.name, err)
		}
		if m == nil {
			t.Fatalf("%s: The returned manager cannot be nil", test.name)
		}
		if len(m.RootPEM()) < 1 {
			t.Errorf("%s: The root PEM can not have a sizeof 0 bytes", test.name)
		}
		if m.RootCrt() == nil {
			t.Errorf("%s: The root cert of the manager cannot be nil", test.name)
		}
		runGenerateCertTest(m, test.name, test.certType, true, t)
		runGenerateCertTest(m, test.name, test.certType, false, t)
		if err := os.RemoveAll(outPath); err != nil {
			t.Fatalf("%s: Cannot remove work dir: %+v", test.name, err)
		}
	}
}
