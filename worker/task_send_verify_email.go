package worker

import (
	"context"
	"encoding/json"
	"fmt"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/BariqDev/ias-bank/util"
	"github.com/hibiken/asynq"
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
	options ...asynq.Option) error {

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

		// if err == pgx.ErrNoRows {
		// 	return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		// }
		return fmt.Errorf("Failed to get user: %w", err)
	}

	// : SEND EMAIL TO USER
	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	// TODO: replace this URL with an environment variable that points to a front-end page
	verifyUrl := fmt.Sprintf("http://localhost:8080/v1/verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("email", user.Email).
		Msg("Task processed")

	return nil

}
