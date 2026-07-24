package worker

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/util"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
}

func (distributor RedisTaskDistributor) DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	log.Info().Str("task_type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Str("task_id", info.ID).
		Msg("task enqueued successfully")
	return nil
}

func (processor *RedisTaskProcessor) ProcessVerifyEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	err := json.Unmarshal(t.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal task payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   payload.Username,
		Email:      payload.Email,
		SecretCode: util.RandomString(32),
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}
	
	subject := "Welcome to Simple Bank! Verify Your Email Address"
	verifyUrl := fmt.Sprintf("http:/localhost:8080/v1//verify_email?email_id=%d&secret_code=%s",
		verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
		Thank your for registering with us!<br/>
		Please<a href="%s">click here</a> to verify your email address.
		`, user.Username, verifyUrl)
	to := []string{payload.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().Str("type", t.Type()).Bytes("payload", t.Payload()).
		Str("email", user.Email).Msg("processed verify email task successfully")
	return nil
}
