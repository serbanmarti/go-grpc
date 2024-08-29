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

// crudUpdateCmd represents the crud-update command
var crudUpdateCmd = &cobra.Command{
	Use:   "crud-update [id] [new_name]",
	Short: "Command to update a resource by ID, changing its name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		runCrudUpdateCmd(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(crudUpdateCmd)
}

func runCrudUpdateCmd(id, newName string) {
	// Create a new client to the CRUD service
	client := internal.NewCrudServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new request for the Update method
	req := connect.NewRequest(&crudv1.UpdateRequest{
		Id:          id,
		UpdatedName: newName,
	})

	// Set the authentication token
	req.Header().Set(environment.TokenHeader, environment.TokenSecret)

	// Call the Update method
	res, err := client.Update(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to update resource: %v\n", err)
	}
	log.Printf("[INFO] Updated resource with ID: %s -> New name: %s\n", res.Msg.Id, res.Msg.Name)
}
