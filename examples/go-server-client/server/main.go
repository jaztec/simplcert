package main

import (
	"context"
	"crypto/tls"
	go_server_client "examples/go-server-client"
	proto "examples/go-server-client/pb"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
)

func main() {
	logrus.Info("Starting server")

	server := grpc.NewServer(grpc.Creds(transportCredentials()))
	defer server.Stop()

	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		logrus.WithField("error", err).Fatal("Error listening")
	}

	proto.RegisterGreeterServiceServer(server, &service{})

	logrus.Fatal(server.Serve(listener))
}

type service struct {
	proto.UnimplementedGreeterServiceServer
}

func (s *service) Greet(_ context.Context, req *proto.GreetingRequest) (*proto.GreetingResponse, error) {
	return &proto.GreetingResponse{
		Greeting: fmt.Sprintf("Hello %s", req.Name),
	}, nil
}

func transportCredentials() credentials.TransportCredentials {
	crt, err := tls.LoadX509KeyPair("/certs/server.crt", "/certs/server.key")
	if err != nil {
		panic(fmt.Errorf("error loading x509 pair: %w", err))
	}
	pool := go_server_client.CAPool()
	credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{crt},
		RootCAs:      pool,
		ServerName:   "Server",
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    pool,
	})
	return credentials.NewServerTLSFromCert(&crt)
}
