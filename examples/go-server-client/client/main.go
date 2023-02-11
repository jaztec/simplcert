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

	pool := go_server_client.CAPool()
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{crt},
		RootCAs:      pool,
	})
}
