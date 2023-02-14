package simplcert_test

import (
	"errors"
	"github.com/jaztec/simplcert"
	"os"
	"testing"
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

		if err := os.RemoveAll(outPath); err != nil {
			t.Fatalf("%s: Cannot remove work dir: %+v", test.name, err)
		}
	}
}
