package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/serbanmarti/go-grpc/client/internal"
	"github.com/serbanmarti/go-grpc/env"
	streamv1 "github.com/serbanmarti/go-grpc/proto_gen/stream/v1"
)

// streamUploadFileCmd represents the stream-upload-file command
var streamUploadFileCmd = &cobra.Command{
	Use:   "stream-upload-file",
	Short: "Command to stream upload a file",
	Run: func(cmd *cobra.Command, args []string) {
		runStreamUploadFileCmd()
	},
}

func init() {
	rootCmd.AddCommand(streamUploadFileCmd)
}

func runStreamUploadFileCmd() {
	// Create a new client to the Stream service
	client := internal.NewStreamServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new stream for the UploadFile method
	stream := client.UploadFile(context.Background())

	// Set the authentication token
	stream.RequestHeader().Set(environment.TokenHeader, environment.TokenSecret)

	// We simulate the file being read in chunks, and requests being created for each chunk
	// In a real-world scenario, the file would be read in chunks from disk
	reqData := []streamv1.UploadFileRequest{
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
	}

	// Send each request to the server
	for _, req := range reqData {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("[ERROR] Failed to stream request: %v\n", err)
		}
	}

	// Close the stream
	res, err := stream.CloseAndReceive()
	if err != nil {
		log.Fatalf("[ERROR] Failed to close and receive from the stream: %v\n", err)
	}
	log.Printf("[INFO] File uploaded! Received confirmation: Filename: %s - Size: %d\n", res.Msg.FileName, res.Msg.Size)
}
