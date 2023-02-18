package simplcert

import (
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

func certConfig(certType CertType, isServer bool) CertConfig {
	return CertConfig{
		Name:     "Test",
		Host:     "test.org, , 127.0.0.1, ::1,",
		IsServer: isServer,
		CertType: certType,
		NotAfter: time.Now().AddDate(0, 0, 1),
	}
}

func TestLoadPublicKey(t *testing.T) {
	tests := []struct {
		name     string
		certType CertType
	}{
		{"Test load ECDSA public key", TypeECDSA},
		{"Test load RSA public key", TypeRSA},
		{"Test load ED25519 public key", TypeED25519},
	}

	for _, test := range tests {
		outPath := createTestDirectory(test.name, t)
		if err := CreateRootCAFiles(test.certType, outPath); err != nil {
			t.Errorf("%s: Creating roto CA files returned an error: %+v", test.name, err)
		}
		m, err := NewManager(outPath)
		if err != nil {
			t.Fatalf("%s: Error creating manager: %+v", test.name, err)
		}

		// create the certificate
		_, priv, pubBytes, err := m.CreateNamedCert(certConfig(test.certType, true))
		if err != nil {
			t.Fatalf("%s: Error creating named certificate: %+v", test.name, err)
		}
		privRaw, err := m.MarshalPrivateKey(priv)
		if err != nil {
			t.Fatalf("%s: Error marshaling private key: %+v", test.name, err)
		}
		privBytes := EncodePrivateKey(privRaw)

		// load the private and public keys
		loadedPriv, err := loadPrivateKey(privBytes)
		if err != nil {
			t.Errorf("%s: Loading the private key file failed: %+v", test.name, err)
		}
		loadedPub, err := loadPublicKey(pubBytes)
		if err != nil {
			t.Errorf("%s: Loading the public key file failed: %+v", test.name, err)
		}
		if !publicKeyEquals(loadedPriv.Public(), loadedPub) {
			t.Errorf("%s: The public keys do not match", test.name)
		}

		if err := os.RemoveAll(outPath); err != nil {
			t.Fatalf("%s: Cannot remove work dir: %+v", test.name, err)
		}
	}
}
