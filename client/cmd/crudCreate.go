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

// crudCreateCmd represents the crud-create command
var crudCreateCmd = &cobra.Command{
	Use:   "crud-create [name]",
	Short: "Command to create a new resource",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runCrudCreateCmd(args[0])
	},
}

func init() {
	rootCmd.AddCommand(crudCreateCmd)
}

func runCrudCreateCmd(name string) {
	// Create a new client to the CRUD service
	client := internal.NewCrudServiceClient()

	// Get the environment configuration
	environment := env.GetEnvironment()

	// Create a new request for the Create method
	req := connect.NewRequest(&crudv1.CreateRequest{
		Name: name,
	})

	// Set the authentication token
	req.Header().Set(environment.TokenHeader, environment.TokenSecret)

	// Call the Create method
	res, err := client.Create(
		context.Background(),
		req,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create resource: %v\n", err)
	}
	log.Printf("[INFO] Created resource with ID: %s\n", res.Msg.Id)
}
