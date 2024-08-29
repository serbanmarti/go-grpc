package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"

	"github.com/serbanmarti/go-grpc/env"
	"github.com/serbanmarti/go-grpc/proto_gen/crud/v1/crudv1connect"
	"github.com/serbanmarti/go-grpc/proto_gen/stream/v1/streamv1connect"
)

func newInsecureClient() *http.Client {
	return &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
		Timeout: 5 * time.Second,
	}
}

func NewCrudServiceClient() crudv1connect.CrudServiceClient {
	// Get the environment configuration
	environment := env.GetEnvironment()

	return crudv1connect.NewCrudServiceClient(
		newInsecureClient(),
		fmt.Sprintf("http://localhost:%d", environment.Port),
		connect.WithGRPC(),
	)
}

func NewStreamServiceClient() streamv1connect.StreamServiceClient {
	// Get the environment configuration
	environment := env.GetEnvironment()

	return streamv1connect.NewStreamServiceClient(
		newInsecureClient(),
		fmt.Sprintf("http://localhost:%d", environment.Port),
		connect.WithGRPC(),
	)
}
