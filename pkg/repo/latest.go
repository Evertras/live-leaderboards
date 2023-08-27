package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	partitionKeyLatest = "latest"

	sortKeyLatestRound = "rnd_id"
)

type Latest struct {
	primaryKey

	ID string `dynamodbav:"id"`
}

func (r *Repo) GetLatestRoundID(ctx context.Context) (string, error) {
	key := primaryKey{
		PK: partitionKeyLatest,
		SK: sortKeyLatestRound,
	}

	avsKey, err := attributevalue.MarshalMap(key)

	if err != nil {
		return "", fmt.Errorf("failed to marshal: %w", err)
	}

	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key:       avsKey,
	})

	if err != nil {
		return "", fmt.Errorf("failed to get item: %w", err)
	}

	var data Latest
	err = attributevalue.UnmarshalMap(result.Item, &data)

	if err != nil {
		return "", fmt.Errorf("failed to unmarshal returned value: %w", err)
	}

	return data.ID, nil
}
