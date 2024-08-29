package interceptor

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"go.uber.org/zap"
)

type LoggerInterceptor struct{}

func NewLoggerInterceptor() *LoggerInterceptor {
	return &LoggerInterceptor{}
}

func (i *LoggerInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(
		ctx context.Context,
		req connect.AnyRequest,
	) (res connect.AnyResponse, err error) {
		// Log the start and end of the request

		// Add the procedure to the log fields
		proc := req.Spec().Procedure
		f := []zap.Field{
			zap.String("procedure", proc),
		}

		zap.L().Info(fmt.Sprintf("started unary request: %s", proc), f...)
		res, err = next(ctx, req)
		if err != nil {
			f = append(f, zap.Error(err), zap.Any("response_code", connect.CodeOf(err)))
		}
		zap.L().Info(fmt.Sprintf("finished unary request: %s", proc), f...)

		return
	}
}

func (i *LoggerInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	// This is a no-op because we don't care about the client side in the server
	return func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		return next(ctx, spec)
	}
}

func (i *LoggerInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) (err error) {
		// Log the start and end of the request

		// Add the procedure to the log fields
		proc := conn.Spec().Procedure
		f := []zap.Field{
			zap.String("procedure", proc),
		}

		zap.L().Info(fmt.Sprintf("started stream request: %s", proc), f...)
		err = next(ctx, conn)
		if err != nil {
			f = append(f, zap.Error(err), zap.Any("response_code", connect.CodeOf(err)))
		}
		zap.L().Info(fmt.Sprintf("finished stream request: %s", proc), f...)

		return
	}
}
