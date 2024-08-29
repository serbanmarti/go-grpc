package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"

	streamv1 "github.com/serbanmarti/go-grpc/proto_gen/stream/v1"
	"github.com/serbanmarti/go-grpc/proto_gen/stream/v1/streamv1connect"
)

func TestStreamService_UploadFile(t *testing.T) {
	tests := []struct {
		name    string
		reqData []streamv1.UploadFileRequest
		resData streamv1.UploadFileResponse
	}{
		{
			name: "Test smaller file",
			reqData: []streamv1.UploadFileRequest{
				{
					FileName: "smaller.txt",
					Chunk:    []byte("Hello, World!"),
				},
				{
					Chunk: []byte("This is a small file."),
				},
				{
					Chunk: []byte("Goodbye!"),
				},
			},
			resData: streamv1.UploadFileResponse{
				FileName: "smaller.txt",
				Size:     42,
			},
		},
		{
			name: "Test larger file",
			reqData: []streamv1.UploadFileRequest{
				{
					Chunk: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."),
				},
				{
					Chunk: []byte("Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat."),
				},
				{
					Chunk: []byte("Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur."),
				},
				{
					FileName: "larger.txt",
					Chunk:    []byte("Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."),
				},
			},
			resData: streamv1.UploadFileResponse{
				FileName: "larger.txt",
				Size:     442,
			},
		},
	}

	client := streamv1connect.NewStreamServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := client.UploadFile(ctx)

			for _, req := range tt.reqData {
				err := stream.Send(&req)
				assert.NoError(t, err)
			}

			res, err := stream.CloseAndReceive()
			assert.NoError(t, err)

			if !cmp.Equal(
				tt.resData, *res.Msg,
				cmpopts.IgnoreUnexported(streamv1.UploadFileResponse{}),
			) {
				t.Errorf("want[-], got[+]\n%v", cmp.Diff(
					tt.resData, *res.Msg,
					cmpopts.IgnoreUnexported(streamv1.UploadFileResponse{}),
				))
			}
		})
	}
}

func TestStreamService_DirectMessage(t *testing.T) {
	tests := []struct {
		name    string
		reqData []streamv1.DirectMessageRequest
		resData []streamv1.DirectMessageResponse
	}{
		{
			name: "Test single message",
			reqData: []streamv1.DirectMessageRequest{
				{
					Message: "Hello",
				},
			},
			resData: []streamv1.DirectMessageResponse{
				{
					Message: "Received message: Hello",
				},
			},
		},
		{
			name: "Test multiple messages",
			reqData: []streamv1.DirectMessageRequest{
				{
					Message: "Hello",
				},
				{
					Message: "How are you?",
				},
			},
			resData: []streamv1.DirectMessageResponse{
				{
					Message: "Received message: Hello",
				},
				{
					Message: "Received message: How are you?",
				},
			},
		},
	}

	client := streamv1connect.NewStreamServiceClient(
		newInsecureClient(),
		"http://127.0.0.1:8080",
		connect.WithGRPC(),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			stream := client.DirectMessage(ctx)

			for idx, req := range tt.reqData {
				err := stream.Send(&req)
				assert.NoError(t, err)

				res, err := stream.Receive()
				assert.NoError(t, err)

				if !cmp.Equal(
					tt.resData[idx], *res,
					cmpopts.IgnoreUnexported(streamv1.DirectMessageResponse{}),
				) {
					t.Errorf("want[-], got[+]\n%v", cmp.Diff(
						tt.resData[idx], *res,
						cmpopts.IgnoreUnexported(streamv1.DirectMessageResponse{}),
					))
				}
			}

			err := stream.CloseRequest()
			assert.NoError(t, err)
			err = stream.CloseResponse()
			assert.NoError(t, err)
		})
	}
}
