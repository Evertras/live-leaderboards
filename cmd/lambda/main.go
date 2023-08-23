package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"

	"github.com/Evertras/live-leaderboards/pkg/repo"
	"github.com/Evertras/live-leaderboards/pkg/server"
)

func main() {
	tableName := os.Getenv("EVERTRAS_LEADERBOARD_TABLE")

	repository, err := repo.NewDefault(context.Background(), tableName)

	if err != nil {
		panic(err)
	}

	server := server.New()

	// Even though we're behind a v2 gateway, we still use the v1 adapter here
	// as we seem to get the v1 event that contains the path, etc.
	lambda.Start(httpadapter.New(server).ProxyWithContext)
}
