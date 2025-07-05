package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/Abelova-Grupa/Mercypher/session-service/external/proto"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/db"
	"github.com/Abelova-Grupa/Mercypher/session-service/internal/grpc/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	// Loading grpc server
	tlsPort := loadGrpcServerPort()
	// creds := loadTransportCredentials()

	listener, err := net.Listen("tcp", tlsPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// grpcServer := grpc.NewServer(grpc.Creds(creds))
	grpcServer := grpc.NewServer()
	pb.RegisterSessionServiceServer(grpcServer, server.NewGrpcServer(db.Connect()))

	go func() {
		log.Printf("Starting gRPC server on port %v...", tlsPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Initialize grpc gateway server for handling handling http requests and converting them to grpc
	httpPort := loadGatewayServerPort()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	clientCreds := loadClientTransportCredentials()

	mux := runtime.NewServeMux()

	err = pb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, "localhost"+tlsPort, []grpc.DialOption{grpc.WithTransportCredentials(clientCreds)})
	if err != nil {
		log.Fatalf("failed to register grpc-gateway: %v", err)
	}

	go func() {
		log.Printf("Starting REST gateway on %s...", httpPort)
		if err := http.ListenAndServe(httpPort, mux); err != nil {
			log.Fatalf("unable to server rest gateway: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	cancel()

}

func loadGrpcServerPort() string {
	tlsPort := ":" + os.Getenv("SESSION_SERVICE_PORT")
	if tlsPort == ":" {
		tlsPort = ":50055"
	}
	return tlsPort
}

func loadTransportCredentials() credentials.TransportCredentials {
	// Loading env variables from cloud
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	certPath := os.Getenv("CERT_PATH")

	var creds credentials.TransportCredentials
	var err error

	// local host execution
	if (tlsCert == "" || tlsKey == "") && certPath == "" {
		creds, err = credentials.NewServerTLSFromFile("internal/certs/localhost.crt", "internal/certs/localhost.key")
		if err != nil {
			log.Fatalf("Failed to load TLS keys: %v", err)
		}
	} else if tlsCert == "" || tlsKey == "" {
		// Docker execution
		creds, err = credentials.NewServerTLSFromFile(certPath+"/localhost.crt", certPath+"/localhost.key")
		if err != nil {
			log.Fatalf("Failed to load TLS keys: %v", err)
		}
		// Railway execution
	} else {
		// Creating a certificate : key pair
		cert, err := tls.X509KeyPair([]byte(tlsCert), []byte(tlsKey))
		if err != nil {
			log.Fatalf("Failed to generate x509 pair: %v", err)
		}
		// Creating tls configuration based on certificate pair
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		// Creating credentials
		creds = credentials.NewTLS(tlsConfig)
	}

	if creds == nil {
		log.Fatal("TLS credentials are nil")
	}

	return creds
}

func loadGatewayServerPort() string {
	httpPort := ":" + os.Getenv("PORT")
	if httpPort == ":" {
		return ":9091"
	}
	return httpPort
}

func loadClientTransportCredentials() credentials.TransportCredentials {
	var creds credentials.TransportCredentials
	var err error
	tslCert := os.Getenv("TLS_CERT")
	certPath := os.Getenv("CERT_PATH")
	if tslCert == "" && certPath == "" {
		creds, err = credentials.NewClientTLSFromFile("internal/certs/localhost.crt", "")
		if err != nil {
			log.Fatalf("unable to create client credentials ")
		}
	} else {
		creds, err = credentials.NewClientTLSFromFile(certPath+"/localhost.crt", "")
		if err != nil {
			log.Fatalf("unable to create client credentials")
		}
		// certPEMBock := []byte(tslCert)
		// certPool := x509.NewCertPool()

		// if ok := certPool.AppendCertsFromPEM(certPEMBock); !ok {
		// 	log.Fatal("failed to append cert to pool")
		// }

		// tlsConfig := &tls.Config{
		// 	RootCAs: certPool,
		// }

		// creds = credentials.NewTLS(tlsConfig)
	}

	return creds
}
