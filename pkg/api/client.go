package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// Some extra methods for the client

func (c *Client) GetRound(ctx context.Context, id uuid.UUID) (*Round, error) {
	response, err := c.GetRoundRoundID(ctx, id.String())

	if err != nil {
		return nil, fmt.Errorf("failed to get round: %w", err)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	raw, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	var round Round

	err = json.Unmarshal(raw, &round)

	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return &round, nil
}

func (c *Client) CreateRound(ctx context.Context, req RoundRequest) (*CreatedRound, error) {
	res, err := c.PostRound(ctx, req)

	if err != nil {
		return nil, fmt.Errorf("t.client.PostRound: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d, expected %d", res.StatusCode, http.StatusCreated)
	}

	raw, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	var createdRound CreatedRound

	err = json.Unmarshal(raw, &createdRound)

	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal result: %w", err)
	}

	return &createdRound, nil
}
