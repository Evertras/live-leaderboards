package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/Evertras/live-leaderboards/pkg/api"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

const (
	sortKeyEventRoundStart = "rnd_start"
)

type EventRoundStart struct {
	RoundID string `dynamodbav:"pk" json:"-"`
	SortKey string `dynamodbav:"sk" json:"-"`

	// TODO: Make this actually detached from the API struct
	// for storage safety... but for now it's nice for simplifying
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
	eventStart, eventScores, err := r.getRoundData(ctx, roundID)

	if err != nil {
		return nil, fmt.Errorf("failed to get round start event: %w", err)
	}

	title := ""

	if eventStart.Title != nil {
		title = *eventStart.Title
	}

	players := make([]api.RoundPlayerData, len(eventStart.Players))

	for i, player := range eventStart.Players {
		players[i].Name = player.Name
	}

	for _, score := range eventScores {
		if score.PlayerIndex < 0 || score.PlayerIndex >= len(players) {
			return nil, fmt.Errorf("found score for player index %d but only have %d players", score.PlayerIndex, len(players))
		}

		players[score.PlayerIndex].Scores = append(players[score.PlayerIndex].Scores, api.HoleScore{
			Hole:  score.HoleNumber,
			Score: score.Score,
		})
	}

	round := api.Round{
		Course:  eventStart.Course,
		Id:      roundID,
		Players: players,
		Title:   title,
	}

	return &round, nil
}

func (r *Repo) getRoundData(ctx context.Context, roundID uuid.UUID) (*EventRoundStart, []EventScore, error) {
	if isZeroUUID(roundID) {
		return nil, nil, fmt.Errorf("id is empty")
	}

	id := roundID.String()

	keyExpression := expression.Key(keyPrimary).Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyExpression).Build()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build query expression: %w", err)
	}

	response, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 r.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to query: %w", err)
	}

	var eventStart EventRoundStart
	var eventScores []EventScore

	for _, item := range response.Items {
		var keyData primaryKey

		err = attributevalue.UnmarshalMap(item, &keyData)

		if err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal key data: %w", err)
		}

		switch {
		case keyData.SK == sortKeyEventRoundStart:
			err = attributevalue.UnmarshalMap(item, &eventStart)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal start event: %w", err)
			}

		case strings.HasPrefix(keyData.SK, sortKeyEventScore):
			var eventScore EventScore
			err = attributevalue.UnmarshalMap(item, &eventScore)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to unmarshal score event: %w", err)
			}
			eventScores = append(eventScores, eventScore)
		}
	}

	return &eventStart, eventScores, nil
}
