package worker

import (
	"context"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/hibiken/asynq"
)

type TaskProcess interface {
	ProcessSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	Start() error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcess {
	server := asynq.NewServer(redisOpt, asynq.Config{})

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
