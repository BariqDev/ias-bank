package worker

import (
	"context"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcess interface {
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func reportError(ctx context.Context, task *asynq.Task, err error) {
	log.Error().Err(err).
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Msg("Process task failed")
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcess {
	server := asynq.NewServer(redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(reportError),
			Logger:       NewLogger(),
		})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessSendVerifyEmail)

	return processor.server.Start(mux)

}
