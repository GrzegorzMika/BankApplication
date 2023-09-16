package worker

import (
	"context"

	"BankApplication/internal/db"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueSendVerificationEmail = "send-verify-email"
	QueueDefault               = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	return &RedisTaskProcessor{
		server: asynq.NewServer(redisOpt, asynq.Config{
			Queues: map[string]int{
				QueueSendVerificationEmail: 10,
				QueueDefault:               5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("task_type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("failed to process task")
			}),
			Logger: NewLogger(),
		}),
		store: store,
	}
}

func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
}
