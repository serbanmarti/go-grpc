package cmd

import (
	"context"
	"log"

	"connectrpc.com/connect"
	"github.com/spf13/cobra"

	"github.com/serbanmarti/go-grpc/client/internal"
	"github.com/serbanmarti/go-grpc/env"
	crudv1 "github.com/serbanmarti/go-grpc/proto_gen/crud/v1"
)

// crudReadCmd represents the crud-read command
var crudReadCmd = &cobra.Command{
	Use:   "crud-read [id]",
	Short: "Command to read a resource by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCrudReadCmd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(crudReadCmd)
}

func runCrudReadCmd(id string) {
	// Create a new client to the CRUD service
	client := internal.NewCrudServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new request for the Read method
	req := connect.NewRequest(&crudv1.ReadRequest{
		Id: id,
	})

	// Set the authentication token
	req.Header().Set(environment.TokenHeader, environment.TokenSecret)

	// Call the Read method
	res, err := client.Read(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read resource: %v\n", err)
	}
	log.Printf("[INFO] Read resource with ID: %s -> Name: %s\n", res.Msg.Id, res.Msg.Name)
}
