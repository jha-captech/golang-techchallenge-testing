package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
	expiration time.Duration
}

func NewClient(client *redis.Client, expiration time.Duration) *Client {
	return &Client{
		Client:     client,
		expiration: expiration,
	}
}

type StringCmd struct {
	*redis.StringCmd
}

type StatusCmd struct {
	*redis.StatusCmd
}

func (c *Client) SetMarshal(ctx context.Context, key string, value interface{}) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("[in services.Client.Set] failed to marshal value: %w", err)
	}

	if err = c.Client.Set(ctx, key, jsonData, c.expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) *StringCmd {
	stringCmd := c.Client.Get(ctx, key)
	return &StringCmd{StringCmd: stringCmd}
}

func (cmd *StringCmd) Result() (string, bool, error) {
	val, err := cmd.StringCmd.Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return "", false, nil
		default:
			return "", false, err
		}
	}

	if val == "" {
		return "", false, nil
	}

	return val, true, nil
}

func (cmd *StringCmd) Unmarshal(v any) (bool, error) {
	val, err := cmd.StringCmd.Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return false, nil
		default:
			return false, err
		}
	}

	if val == "" {
		return false, nil
	}

	if err = json.Unmarshal([]byte(val), v); err != nil {
		return false, fmt.Errorf("[in services.StringCmd.Unmarshal] failed to unmarshal from cache: %w", err)
	}

	return true, nil
}
