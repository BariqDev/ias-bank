package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributer interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		options ...asynq.Option,) error
}

type RedisTaskDistributer struct {
	client *asynq.Client
}

func NewRedisTaskDistributer(redisOpt asynq.RedisClientOpt) TaskDistributer {

	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributer{
		client: client,
	}
}
