package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "simpletodo/db/sqlc"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerifyEmail struct {
	UserName string `json:"user_name"`
}

func (distributor *RedisTaskDistributor)DistributeTaskSendVerifyEmail(
	ctx context.Context,
	payload *PayloadSendVerifyEmail,
	opts ...asynq.Option,
)error{

	jasoPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("faild to marshal task payload: %w", err)
	}
	task := asynq.NewTask(TaskSendVerifyEmail, jasoPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx,task)
	if err != nil {
		return fmt.Errorf("faild to enqueue task: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
	Str("queue", info.Queue).Int("max_retry", info.MaxRetry).Msg("enqueued task")
	return nil
}


func (processor *RedisTaskProcessor)ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error{
	var payload PayloadSendVerifyEmail
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("faild to unmarshal payload: %w",asynq.SkipRetry)
	}

	fmt.Println("Payload:", payload)

    // Ensure GetAccountsParams is properly defined in the db package
    params := db.GetAccountsParams{
        UserName: payload.UserName,
    }

    // Debugging statement to ensure params is correctly constructed
    fmt.Println("Params:", params)

	user, err := processor.store.GetAccount(ctx, params)
	if err != nil {
		if err == sql.ErrNoRows{
			return fmt.Errorf("user doesn't exits: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("faild to get account: %w", err)
	}

	//email to user
	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
	Str("email", user.Account.Email).Msg("processed task")

	return nil
}
