package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Repo struct {
	client    *dynamodb.Client
	tableName *string
}

func NewDefault(ctx context.Context, tableName string, opts ...func(o *dynamodb.Options)) (*Repo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg, opts...)

	return NewWithClient(client, tableName), nil
}

func NewWithClient(client *dynamodb.Client, tableName string) *Repo {
	return &Repo{
		client:    client,
		tableName: &tableName,
	}
}
