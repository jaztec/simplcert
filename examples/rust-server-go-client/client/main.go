package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	proto "examples/rust-server-go-client/client/pb"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"time"
)

func main() {
	logrus.Info("Starting client")

	conn, err := grpc.Dial("server:8000", grpc.WithTransportCredentials(transportCredentials()))
	if err != nil {
		panic(err)
	}
	client := proto.NewGreeterServiceClient(conn)

	fn := func(client proto.GreeterServiceClient) {
		resp, err := client.Greet(context.Background(), &proto.GreetingRequest{Name: "World"})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Greeter: %s\n", resp.Greeting)
	}

	fn(client)

	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			fn(client)
		}
	}
}

func transportCredentials() credentials.TransportCredentials {
	crt, err := tls.LoadX509KeyPair("/certs/client.crt", "/certs/client.key")
	if err != nil {
		panic(fmt.Errorf("error loading x509 pair: %w", err))
	}

	pool := caPool()
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{crt},
		RootCAs:      pool,
	})
}

func caPool() *x509.CertPool {
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
