package repo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"

	"github.com/Evertras/live-leaderboards/pkg/api"
)

const (
	sortKeyEventScore = "score"
)

type EventScore struct {
	RoundID string `dynamodbav:"pk"`
	SortKey string `dynamodbav:"sk"`

	HoleNumber  int `dynamodbav:"h"`
	PlayerIndex int `dynamodbav:"i"`
	Score       int `dynamodbav:"s"`
}

func (r *Repo) SetScore(ctx context.Context, roundID uuid.UUID, req api.PlayerScoreEvent) error {
	data := EventScore{
		RoundID:     roundID.String(),
		SortKey:     sortKeyEventScore,
		HoleNumber:  req.Hole,
		PlayerIndex: req.PlayerIndex,
		Score:       req.Score,
	}

	avs, err := attributevalue.MarshalMap(data)

	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	_, err = r.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			Item:      avs,
			TableName: r.tableName,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}
