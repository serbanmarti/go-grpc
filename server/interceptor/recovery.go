package interceptor

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type RecoveryInterceptor struct{}

func NewRecoveryInterceptor() *RecoveryInterceptor {
	return &RecoveryInterceptor{}
}

func (i *RecoveryInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(
		ctx context.Context,
		req connect.AnyRequest,
	) (res connect.AnyResponse, err error) {
		defer func() {
			// Recover from any panics and return an internal server error
			if r := recover(); r != nil {
				zap.L().Error("recovered from panic", zap.Any("panic", r))
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("unexpected server error"))
			}
		}()
		return next(ctx, req)
	}
}

func (i *RecoveryInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	// This is a no-op because we don't care about the client side in the server
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		return next(ctx, spec)
	}
}

func (i *RecoveryInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) (err error) {
		defer func() {
			// Recover from any panics and return an internal server error
			if r := recover(); r != nil {
				zap.L().Error("recovered from panic", zap.Any("panic", r))
				err = connect.NewError(connect.CodeInternal, fmt.Errorf("unexpected server error"))
			}
		}()
		return next(ctx, conn)
	}
}
