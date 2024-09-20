package repository

import (
	"context"
)

type WebhookThread struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	ThreadID        string `db:"thread_id"`
}

func (r *Repository) InsertWebhookThread(
	ctx context.Context,
	webhookSerialID int64,
	threadID string,
) error {
	query := `
		INSERT INTO webhook_thread (
			webhook_serial_id,
			thread_id
		) VALUES (
			$1,
			$2
		) ON CONFLICT (webhook_serial_id, thread_id) DO NOTHING
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		webhookSerialID,
		threadID,
	)
	return err
}

func (r *Repository) GetWebhookThreadWithWebhookSerialID(
	ctx context.Context,
	webhookSerialID int64,
) ([]*WebhookThread, error) {
	query := `
		SELECT
			*
		FROM
			webhook_thread
		WHERE
			webhook_serial_id = $1
	`
	var webhookThread []*WebhookThread
	err := r.db.SelectContext(ctx, &webhookThread, query, webhookSerialID)
	return webhookThread, err
}

func (r *Repository) DeleteWebhookThreadsNotInProvidedList(
	ctx context.Context,
	webhookSerialID int64,
) error {
	query := `
		DELETE FROM
			webhook_thread
		WHERE
			webhook_serial_id = $1
	`
	_, err := r.db.ExecContext(ctx, query, webhookSerialID)
	return err
}
