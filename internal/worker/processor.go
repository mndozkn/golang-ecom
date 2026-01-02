package worker

import (
	"context"
	"encoding/json"
	"go-crud/internal/infrastructure/mailer"
	"log"

	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	server *asynq.Server
	mailer mailer.Mailer
}

func NewTaskProcessor(redisOpt asynq.RedisClientOpt, mailer mailer.Mailer) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{Concurrency: 5})
	return &TaskProcessor{server: server, mailer: mailer}
}

func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc("task:send_welcome_email", p.processSendWelcomeEmail)
	return p.server.Run(mux)
}

func (p *TaskProcessor) processSendWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var payload map[string]string
	json.Unmarshal(t.Payload(), &payload)

	log.Printf("📧 Mail gönderiliyor: %s", payload["email"])
	return p.mailer.SendWelcomeEmail(payload["email"])
}
