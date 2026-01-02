package worker

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendWelcomeEmail(ctx context.Context, email string) error
}

type redisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	return &redisTaskDistributor{
		client: asynq.NewClient(redisOpt),
	}
}

func (d *redisTaskDistributor) DistributeTaskSendWelcomeEmail(ctx context.Context, email string) error {
	payload, _ := json.Marshal(map[string]string{"email": email})
	task := asynq.NewTask("task:send_welcome_email", payload)
	_, err := d.client.EnqueueContext(ctx, task)
	return err
}
