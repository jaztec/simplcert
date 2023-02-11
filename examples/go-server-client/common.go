package go_server_client

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func CAPool() *x509.CertPool {
	path := "/certs/root-ca.crt"
	r, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(r)
	if block == nil {
		panic(fmt.Errorf("decoding file from %s failed", path))
	}
	crt, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(crt)

	return pool
}
