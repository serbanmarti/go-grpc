package service

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/serbanmarti/go-grpc/proto_gen/crud/v1/crudv1connect"
	"github.com/serbanmarti/go-grpc/proto_gen/stream/v1/streamv1connect"
)

func init() {
	// Create the mock data store
	data := make(map[string]string)
	data["2imgNBCejbjXehOazVerssNsgcz"] = "Test Record 1"
	data["2imgN7lkpYjE16akMMn52Uvkgln"] = "Test Record 2"

	// Create the server mux & register the services we want to test
	mux := http.NewServeMux()
	mux.Handle(crudv1connect.NewCrudServiceHandler(&CrudService{
		Data:  data,
		Mutex: sync.RWMutex{},
	}))
	mux.Handle(streamv1connect.NewStreamServiceHandler(&StreamService{}))

	go func() {
		http.ListenAndServe(
			"0.0.0.0:8080",
			// Use h2c so we can serve HTTP/2 without TLS.
			h2c.NewHandler(mux, &http2.Server{}),
		)
	}()
}

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
