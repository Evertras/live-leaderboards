package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/Evertras/live-leaderboards/pkg/repo"
	"github.com/Evertras/live-leaderboards/pkg/server"
)

func main() {
	tablename := os.Getenv("EVERTRAS_LEADERBOARD_TABLE")

	r, err := repo.NewDefault(context.Background(), tablename, func(o *dynamodb.Options) {
		url := "http://dynamodb-local:8000"
		o.BaseEndpoint = &url
	})

	if err != nil {
		panic(err)
	}

	s := server.New(r)

	log.Fatal(s.ListenAndServe(":8037"))
}
