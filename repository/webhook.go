package repository

import (
	"context"
	"time"
)

type Webhook struct {
	ID               *int64    `db:"id"`
	GuildID          string    `db:"guild_id"`
	WebhookID        string    `db:"webhook_id"`
	SubscriptionType string    `db:"subscription_type"`
	SubscriptionID   string    `db:"subscription_id"`
	LastPostedAt     time.Time `db:"last_posted_at"`
}

func NewWebhook(
	id *int64,
	guildID string,
	webhookID string,
	subscriptionType string,
	subscriptionID string,
	lastPostedAt time.Time,
) *Webhook {
	return &Webhook{
		ID:               id,
		GuildID:          guildID,
		WebhookID:        webhookID,
		SubscriptionType: subscriptionType,
		SubscriptionID:   subscriptionID,
		LastPostedAt:     lastPostedAt,
	}
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
	id int64,
	lastPostedAt time.Time,
) error {
	query := `
		UPDATE
			webhook
		SET
			last_posted_at = $1
		WHERE
			id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		lastPostedAt,
		id,
	)
	return err
}

func (r *Repository) UpdateWebhookWithWebhookIDAndSubscription(
	ctx context.Context,
	id int64,
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
			id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		subscriptionID,
		subscriptionType,
		id,
	)
	return err
}

func (r *Repository) DeleteWebhook(
	ctx context.Context,
	id int64,
) error {
	query := `
		DELETE FROM
			webhook
		WHERE
			id = $1
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)
	return err
}
