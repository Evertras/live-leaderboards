package repo

import (
	"context"
	"fmt"

	"github.com/Evertras/live-leaderboards/pkg/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

const (
	sortKeyEventRoundStart = "rnd_start"
)

type primaryKey struct {
	PK string `dynamodbav:"pk" json:"-"`
	SK string `dynamodbav:"sk" json:"-"`
}

type EventRoundStart struct {
	RoundID string `dynamodbav:"pk" json:"-"`
	SortKey string `dynamodbav:"sk" json:"-"`

	api.RoundRequest
}

func (r *Repo) CreateEventRoundStart(ctx context.Context, roundID uuid.UUID, req api.RoundRequest) error {
	data := EventRoundStart{
		RoundID:      roundID.String(),
		SortKey:      sortKeyEventRoundStart,
		RoundRequest: req,
	}

	avs, err := attributevalue.MarshalMap(data)

	if err != nil {
		return fmt.Errorf("attributevalue.MarshalMap: %w", err)
	}

	_, err = r.client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			Item:      avs,
			TableName: r.tableName,
		},
	)

	if err != nil {
		return fmt.Errorf("r.client.PutItem: %w", err)
	}

	return nil
}

func (r *Repo) GetRound(ctx context.Context, roundID uuid.UUID) (*api.Round, error) {
	eventStart, err := r.getEventRoundStart(ctx, roundID)

	if err != nil {
		fmt.Errorf("failed to get round start event: %w", err)
	}

	title := ""

	if eventStart.Title != nil {
		title = *eventStart.Title
	}

	players := make([]api.RoundPlayerData, len(eventStart.Players))

	for i, player := range eventStart.Players {
		players[i].Name = player.Name

		// TODO: Add scores here
	}

	round := api.Round{
		Course:  eventStart.Course,
		Id:      roundID,
		Players: players,
		Title:   title,
	}

	return &round, nil
}

func (r *Repo) getEventRoundStart(ctx context.Context, roundID uuid.UUID) (*EventRoundStart, error) {
	key := primaryKey{
		PK: roundID.String(),
		SK: sortKeyEventRoundStart,
	}

	avs, err := attributevalue.MarshalMap(key)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		Key:       avs,
		TableName: r.tableName,
	})

	if err != nil {
		return nil, fmt.Errorf("r.client.GetItem: %w", err)
	}

	var data EventRoundStart

	err = attributevalue.UnmarshalMap(out.Item, &data)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return &data, nil
}
