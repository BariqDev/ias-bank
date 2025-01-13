package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// create distributer
const (
	TaskSendVerifyEmail = "task:send-verify-email"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (distributer *RedisTaskDistributer) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	options ...asynq.Option,) error {

	jsonPayloadBytes, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayloadBytes, options...)

	info, err := distributer.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("Failed to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", jsonPayloadBytes).
		Str("Queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("Task enqueued")

	return nil
}

func (processor *RedisTaskProcessor) ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error {

	var payload PayloadSendVerifyEmail

	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {

		if err == pgx.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("Failed to get user: %w", err)
	}

	// TODO: SEND EMAIL TO USER

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("Task processed")

	return nil

}
