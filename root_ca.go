package simplcert

import (
	"crypto"
	"crypto/x509"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	privateCAKeyFile = rootCABaseFilename + ".key"
	publicCAKeyFile  = rootCABaseFilename + ".pub"
	certCAKeyFile    = rootCABaseFilename + ".crt"
)

var baseIdentifier = []int{1, 0, 9652}

func loadCA(certsFolder string) (*x509.Certificate, crypto.Signer, error) {
	if !checkRootCAFiles(certsFolder) {
		return nil, nil, fmt.Errorf("certificates not found in %s", certsFolder)
	}

	crt, err := loadRootCACertificate(certsFolder)
	if err != nil {
		return nil, nil, err
	}

	key, err := loadRootCAPrivateKey(certsFolder)
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func loadRootCACertificate(certsFolder string) (*x509.Certificate, error) {
	crtPath := certsFolder + string(os.PathSeparator) + certCAKeyFile

	b, err := os.ReadFile(crtPath)
	if err != nil {
		return nil, err
	}

	crt, err := loadCert(b)
	if err != nil {
		return nil, err
	}

	return crt, nil
}

func loadRootCAPrivateKey(certsFolder string) (crypto.Signer, error) {
	keyPath := certsFolder + string(os.PathSeparator) + privateCAKeyFile

	b, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return loadPrivateKey(b)
}

func checkRootCAFiles(outPath string) bool {
	files := []string{
		outPath + string(os.PathSeparator) + privateCAKeyFile,
		outPath + string(os.PathSeparator) + publicCAKeyFile,
		outPath + string(os.PathSeparator) + certCAKeyFile,
	}
	for _, f := range files {
		if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}

func CreateRootCAFiles(certType CertType, outPath string) (err error) {
	log.WithField("out_path", outPath).Info("Creating root CA files")
	crt, priv, err := createRootCertificate(certType)
	if err != nil {
		return
	}

	crtPem := EncodeCertificate(crt.Raw)
	if err = os.WriteFile(outPath+string(os.PathSeparator)+certCAKeyFile, crtPem, 0644); err != nil {
		return
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return
	}
	key := EncodePrivateKey(privBytes)
	if err = os.WriteFile(outPath+string(os.PathSeparator)+privateCAKeyFile, key, 0644); err != nil {
		return
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return
	}
	pub := EncodePublicKey(pubBytes)
	if err = os.WriteFile(outPath+string(os.PathSeparator)+publicCAKeyFile, pub, 0644); err != nil {
		return
	}

	return
}

func createRootCertificate(certType CertType) (*x509.Certificate, crypto.Signer, error) {
	priv, err := createKey(certType)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate key")
	}

	crt, err := createNamedCert(CertConfig{
		Name:     "Root CA",
		usage:    x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		extUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:     true,
		NotAfter: time.Now().AddDate(1, 0, 0),
	}, nil, priv.Public(), priv, baseIdentifier)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate root key")
	}

	return crt, priv, nil
}
