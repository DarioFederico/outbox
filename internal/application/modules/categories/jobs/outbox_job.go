package jobs

import (
	"context"
	"database/sql"
	"time"

	"outbox/config"
	"outbox/internal/application/modules/categories/models"
	"outbox/internal/infrastructure/log"
	"outbox/internal/infrastructure/mbaas"
)

const (
	getMessageToProcessList string = `SELECT id, type, message, status, created_at, updated_at FROM outbox WHERE status = ? LIMIT 10 FOR UPDATE`
	updateMessage           string = `UPDATE outbox SET status = ?, updated_at = ? WHERE id = ?`
)

type outboxJob struct {
	cfg    *config.AppConfig
	client mbaas.MBaaS
	db     *sql.DB
}

func NewOutboxJob(cfg *config.AppConfig, client mbaas.MBaaS, db *sql.DB) Job {
	return &outboxJob{cfg: cfg, client: client, db: db}
}

func (j *outboxJob) Run() {
	interval, _ := time.ParseDuration(j.cfg.OutboxInterval)
	go func() {
		for {
			ctx := context.Background()
			if err := j.process(ctx); err != nil {
				log.For(ctx).Errorf("error running job. %+v", err)
			}
			time.Sleep(interval)
		}
	}()
}

func (j *outboxJob) process(ctx context.Context) error {
	tx, err := j.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	messages, err := j.getMessagesToProcess(ctx, tx)
	if err != nil {
		log.For(ctx).Errorf("error getting messages to process. %+v", err)
		return err
	}

	if len(messages) == 0 {
		log.For(ctx).Info("no messages to process")
		return nil
	}

	for _, message := range messages {
		log.For(ctx).Infof("message to send: %+v", message)
		if err != nil {
			log.For(ctx).Errorf("error updating messages. %+v", err)
			continue
		}

		err = j.client.Publish(ctx, message)
		if err != nil {
			log.For(ctx).Errorf("error publish messages to MBaaS. %+v", err)
			continue
		}

		log.For(ctx).Infof("message was send. %+v", message)
		message.Status = models.OutboxStatus_Processed
		message.UpdatedAt = time.Now()
		err = j.updateMessage(ctx, tx, message)
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (j *outboxJob) getMessagesToProcess(ctx context.Context, tx *sql.Tx) ([]models.Outbox, error) {
	rows, err := tx.Query(getMessageToProcessList, models.OutboxStatus_Pending)
	if err != nil {
		log.For(ctx).Errorf("error to get outbox messages. %+v", err)
		return nil, err
	}

	defer rows.Close()

	var messages []models.Outbox
	var outbox models.Outbox
	for rows.Next() {
		err = rows.Scan(
			&outbox.ID,
			&outbox.Type,
			&outbox.Message,
			&outbox.Status,
			&outbox.CreatedAt,
			&outbox.UpdatedAt,
		)
		if err != nil {
			log.For(ctx).Errorf("error to scan outbox message. %+v", err)
			return messages, err
		}
		messages = append(messages, outbox)
	}

	return messages, nil
}

func (j *outboxJob) updateMessage(ctx context.Context, tx *sql.Tx, outbox models.Outbox) error {
	_, err := tx.Exec(updateMessage,
		outbox.Status,
		outbox.UpdatedAt,
		outbox.ID,
	)
	if err != nil {
		log.For(ctx).Errorf("error update outbox message %+v", err)
		return err
	}
	return nil
}
