package worker

import (
	"context"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/mail"
	"github.com/devder/grpc-b/util"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
	config util.Config
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Error().
					Err(err).
					Str("type", task.Type()).
					Bytes("payload", task.Payload()).
					Msg("process task failed")
			}),
			Logger: NewLogger(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
		config: config,
		mailer: mailer,
	}
}

func (p *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(taskSendVerifyEmail, p.ProcessTaskSendVerifyEmail)

	return p.server.Start(mux)
}
