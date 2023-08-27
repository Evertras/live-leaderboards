package repo_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ory/dockertest"

	"github.com/Evertras/live-leaderboards/pkg/repo"
)

const (
	testTableName = "testtable"
)

var (
	dynamoDBPort string
)

func ptr[K any](item K) *K {
	return &item
}

func TestMain(m *testing.M) {
	// Overwrite any env vars to be safe, and also avoid needing to do that
	// when running tests in CI, etc anyway
	os.Setenv("AWS_ACCESS_KEY_ID", "DUMMYID")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "DUMMYKEY")
	os.Setenv("AWS_DEFAULT_REGION", "us-east-1")

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Failed to start docker pool: %v", err)
	}

	res, err := pool.Run("amazon/dynamodb-local", "1.22.0", nil)
	if err != nil {
		log.Fatalf("Failed to run dynamodb-local: %v", err)
	}

	err = res.Expire(500)
	if err != nil {
		log.Fatalf("Failed to set expire: %v", err)
	}

	dynamoDBPort = res.GetPort("8000/tcp")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)

	client, err := newDynamoDBClient(ctx)

	if err != nil {
		log.Fatalf("Failed to get DynamoDB client")
	}

	pool.MaxWait = time.Minute

	err = pool.Retry(func() error {
		_, err := client.ListTables(ctx, &dynamodb.ListTablesInput{})
		return err
	})

	if err != nil {
		log.Fatalf("Failed to list tables: %v", err)
	}

	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("pk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("sk"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("pk"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("sk"),
				KeyType:       types.KeyTypeRange,
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
	})
	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}

	code := m.Run()

	if err := pool.Purge(res); err != nil {
		log.Fatalf("Failed to purge resource: %v", err)
	}
	cancel()

	os.Exit(code)
}

func newDynamoDBClient(ctx context.Context) (*dynamodb.Client, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           "http://localhost:" + dynamoDBPort,
				SigningRegion: "us-east-1",
			}, nil
		},
	)
	cfg, err := config.LoadDefaultConfig(ctx, config.WithEndpointResolverWithOptions(resolver))

	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}
	client := dynamodb.NewFromConfig(cfg)

	return client, nil
}

func newRepo(ctx context.Context, t *testing.T) *repo.Repo {
	t.Helper()

	client, err := newDynamoDBClient(ctx)

	if err != nil {
		t.Fatalf("Failed to get DynamoDB client: %v", err)
	}

	return repo.NewWithClient(client, testTableName)
}
