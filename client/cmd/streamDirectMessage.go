package cmd

import (
	"context"
	"log"

	"github.com/spf13/cobra"

	"github.com/serbanmarti/go-grpc/client/internal"
	"github.com/serbanmarti/go-grpc/env"
	streamv1 "github.com/serbanmarti/go-grpc/proto_gen/stream/v1"
)

// streamDirectMessageCmd represents the stream-direct-message command
var streamDirectMessageCmd = &cobra.Command{
	Use:   "stream-direct-message",
	Short: "Command to stream a direct message",
	Run: func(cmd *cobra.Command, args []string) {
		runStreamDirectMessageCmd()
	},
}

func init() {
	rootCmd.AddCommand(streamDirectMessageCmd)
}

func runStreamDirectMessageCmd() {
	// Create a new client to the Stream service
	client := internal.NewStreamServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new stream for the DirectMessage method
	stream := client.DirectMessage(context.Background())

	// Set the authentication token
	stream.RequestHeader().Set(environment.TokenHeader, environment.TokenSecret)

	// We simulate the messages being read from input, and requests being created for each message
	// In a real-world scenario, the messages would be read from a file, or a database, or any other source
	reqData := []streamv1.DirectMessageRequest{
		{
			Message: "Hello",
		},
		{
			Message: "How are you?",
		},
	}

	// Send each request to the server and receive the response
	for _, req := range reqData {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalf("[ERROR] Failed to stream direct message: %v\n", err)
		}

		res, err := stream.Receive()
		if err != nil {
			log.Fatalf("[ERROR] Failed to receive direct message response: %v\n", err)
		}
		log.Printf("[INFO] Received direct message response: <%s>\n", res.Message)
	}

	// Close the stream
	err := stream.CloseRequest()
	if err != nil {
		log.Fatalf("[ERROR] Failed to close the stream request: %v\n", err)
	}
	err = stream.CloseResponse()
	if err != nil {
		log.Fatalf("[ERROR] Failed to close the stream reasponse: %v\n", err)
	}
	log.Printf("[INFO] Stream closed successfully\n")
}
