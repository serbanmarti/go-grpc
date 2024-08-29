package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"connectrpc.com/connect"
	"go.uber.org/zap"

	streamv1 "github.com/serbanmarti/go-grpc/proto_gen/stream/v1"
)

type StreamService struct{}

func (s *StreamService) UploadFile(ctx context.Context, stream *connect.ClientStream[streamv1.UploadFileRequest]) (*connect.Response[streamv1.UploadFileResponse], error) {
	// Initialize variables for data we care about from the stream
	fileName := ""
	fileSize := 0

	// Receive data from client
	for stream.Receive() {
		// Use only the first file name received
		if fileName == "" && stream.Msg().GetFileName() != "" {
			fileName = stream.Msg().GetFileName()
		}

		// Calculate the total size of the file
		chunk := stream.Msg().GetChunk()
		fileSize += len(chunk)
	}

	// Check for any errors during the stream
	if err := stream.Err(); err != nil {
		zap.L().Error("Error receiving stream", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("error receiving stream"))
	}

	return connect.NewResponse(&streamv1.UploadFileResponse{
		FileName: fileName,
		Size:     uint32(fileSize),
	}), nil
}

func (s *StreamService) DirectMessage(ctx context.Context, stream *connect.BidiStream[streamv1.DirectMessageRequest, streamv1.DirectMessageResponse]) error {
	for {
		// Receive data from client
		req, err := stream.Receive()
		if errors.Is(err, io.EOF) { // Client closed stream, we need to return nil to indicate success
			return nil
		}
		if err != nil {
			zap.L().Error("Error receiving stream", zap.Error(err))
			return connect.NewError(connect.CodeInternal, fmt.Errorf("error receiving stream"))
		}

		// Send data to client
		err = stream.Send(&streamv1.DirectMessageResponse{
			Message: fmt.Sprintf("Received message: %s", req.GetMessage()),
		})
		if err != nil {
			zap.L().Error("Error sending stream", zap.Error(err))
			return connect.NewError(connect.CodeInternal, fmt.Errorf("error sending stream"))
		}
	}
}
