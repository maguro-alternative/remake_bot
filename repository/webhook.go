package repository

import (
	"context"
	"time"
)

type Webhook struct {
	WebhookSerialID  *int64    `db:"webhook_serial_id"`
	GuildID          string    `db:"guild_id"`
	WebhookID        string    `db:"webhook_id"`
	SubscriptionType string    `db:"subscription_type"`
	SubscriptionID   string    `db:"subscription_id"`
	LastPostedAt     time.Time `db:"last_posted_at"`
}

func (r *Repository) InsertWebhook(
	ctx context.Context,
	guildID string,
	webhookID string,
	subscriptionType string,
	subscriptionID string,
	lastPostedAt time.Time,
) error {
	query := `
		INSERT INTO webhook (
			guild_id,
			webhook_id,
			subscription_type,
			subscription_id,
			last_posted_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		guildID,
		webhookID,
		subscriptionType,
		subscriptionID,
		lastPostedAt,
	)
	return err
}

func (r *Repository) GetAllColumnsWebhooksByGuildID(
	ctx context.Context,
	guildID string,
) ([]*Webhook, error) {
	query := `
		SELECT
			*
		FROM
			webhook
		WHERE
			guild_id = $1
	`
	var webhooks []*Webhook
	err := r.db.SelectContext(
		ctx,
		&webhooks,
		query,
		guildID,
	)
	return webhooks, err
}

func (r *Repository) UpdateWebhookWithLastPostedAt(
	ctx context.Context,
	webhookSerialID int64,
	lastPostedAt time.Time,
) error {
	query := `
		UPDATE
			webhook
		SET
			last_posted_at = $1
		WHERE
			webhook_serial_id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		lastPostedAt,
		webhookSerialID,
	)
	return err
}

func (r *Repository) UpdateWebhookWithWebhookIDAndSubscription(
	ctx context.Context,
	webhookSerialID int64,
	webhookID string,
	subscriptionID string,
	subscriptionType string,
) error {
	query := `
		UPDATE
			webhook
		SET
			webhook_id = $1,
			subscription_id = $2,
			subscription_type = $3
		WHERE
			webhook_serial_id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		subscriptionID,
		subscriptionType,
		webhookSerialID,
	)
	return err
}

func (r *Repository) DeleteWebhookByWebhookSerialID(
	ctx context.Context,
	webhookSerialID int64,
) error {
	query := `
		DELETE FROM
			webhook
		WHERE
			webhook_serial_id = $1
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		webhookSerialID,
	)
	return err
}
