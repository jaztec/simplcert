package manager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

const (
	privateCAKeyFile = rootCABaseFilename + ".key"
	publicCAKeyFile  = rootCABaseFilename + ".pub"
	certCAKeyFile    = rootCABaseFilename + ".crt"
)

func loadCA(certsFolder string) (*x509.Certificate, *ecdsa.PrivateKey, []byte, error) {
	if !checkRootCAFiles(certsFolder) {
		return nil, nil, nil, fmt.Errorf("certificates not found in %s", certsFolder)
	}

	crt, crtPem, err := loadRootCACertificate(certsFolder)
	if err != nil {
		return nil, nil, nil, err
	}

	key, err := loadRootCAPrivateKey(certsFolder)
	if err != nil {
		return nil, nil, nil, err
	}

	return crt, key, crtPem, nil
}

func loadRootCACertificate(certsFolder string) (*x509.Certificate, []byte, error) {
	crtPath := certsFolder + string(os.PathSeparator) + certCAKeyFile

	b, err := os.ReadFile(crtPath)
	if err != nil {
		return nil, nil, err
	}

	crt, err := loadCert(b)
	if err != nil {
		return nil, nil, err
	}

	return crt, b, nil
}

func loadRootCAPrivateKey(certsFolder string) (*ecdsa.PrivateKey, error) {
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

func createRootCAFiles(outPath string) (err error) {
	log.WithField("out_path", outPath).Info("Creating root CA files")
	_, crtPem, priv, err := createRootCertificate()
	if err != nil {
		return
	}
	if err = os.WriteFile(outPath+string(os.PathSeparator)+certCAKeyFile, crtPem, 0644); err != nil {
		return
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return
	}
	key := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privBytes,
		},
	)
	if err = os.WriteFile(outPath+string(os.PathSeparator)+privateCAKeyFile, key, 0644); err != nil {
		return
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return
	}
	pub := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		},
	)
	if err = os.WriteFile(outPath+string(os.PathSeparator)+publicCAKeyFile, pub, 0644); err != nil {
		return
	}

	return
}

func createRootCertificate() (*x509.Certificate, []byte, *ecdsa.PrivateKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate key")
	}

	crt, crtPem, err := createNamedCert(CertConfig{
		Name:     "Root CA",
		usage:    x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		extUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:     true,
	}, nil, &priv.PublicKey, priv)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate root key")
	}
	return crt, crtPem, priv, nil
}
