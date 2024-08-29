package interceptor

import (
	"context"
	"fmt"

	"connectrpc.com/connect"

	"github.com/serbanmarti/go-grpc/env"
)

var errNoToken = fmt.Errorf("auth token missing or invalid")

type AuthInterceptor struct {
	secret string
	header string
}

func NewAuthInterceptor() *AuthInterceptor {
	environment := env.GetEnvironment()
	return &AuthInterceptor{
		secret: environment.TokenSecret,
		header: environment.TokenHeader,
	}
}

func (i *AuthInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		reqToken := req.Header().Get(i.header)
		if reqToken != i.secret {
			return nil, connect.NewError(connect.CodeUnauthenticated, errNoToken)
		}
		return next(ctx, req)
	}
}

func (i *AuthInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	// This is a no-op because we don't care about the client side in the server
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		return next(ctx, spec)
	}
}

func (i *AuthInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		reqToken := conn.RequestHeader().Get(i.header)
		if reqToken != i.secret {
			return connect.NewError(connect.CodeUnauthenticated, errNoToken)
		}
		return next(ctx, conn)
	}
}
