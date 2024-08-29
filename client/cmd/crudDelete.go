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

// crudDeleteCmd represents the crud-delete command
var crudDeleteCmd = &cobra.Command{
	Use:   "crud-delete [id]",
	Short: "Command to delete a resource by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCrudDeleteCmd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(crudDeleteCmd)
}

func runCrudDeleteCmd(id string) {
	// Create a new client to the CRUD service
	client := internal.NewCrudServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new request for the Delete method
	req := connect.NewRequest(&crudv1.DeleteRequest{
		Id: id,
	})

	// Set the authentication token
	req.Header().Set(environment.TokenHeader, environment.TokenSecret)

	// Call the Delete method
	res, err := client.Delete(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to delete resource: %v\n", err)
	}
	log.Printf("[INFO] Deleted resource with ID: %s\n", res.Msg.Id)
}
