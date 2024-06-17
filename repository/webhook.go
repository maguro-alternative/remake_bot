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

func (r *Repository) InsertWebhook(ctx context.Context, webhook *Webhook) error {
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
		webhook.GuildID,
		webhook.WebhookID,
		webhook.SubscriptionType,
		webhook.SubscriptionID,
		webhook.LastPostedAt,
	)
	return err
}
