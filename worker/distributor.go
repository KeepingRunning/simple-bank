package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistibutor interface{
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}


func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistibutor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
