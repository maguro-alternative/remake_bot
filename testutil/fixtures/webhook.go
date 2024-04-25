package fixtures

import (
	"context"
	"testing"
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

func NewWebhook(ctx context.Context, setter ...func(b *Webhook)) *ModelConnector {
	webhook := &Webhook{
		GuildID:          "1111111111111",
		WebhookID:        "1111111111111",
		SubscriptionType: "niconico",
		SubscriptionID:   "1111111111111",
		LastPostedAt:     time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
	}

	return &ModelConnector{
		Model: webhook,
		setter: func() {
			for _, s := range setter {
				s(webhook)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.Webhooks = append(f.Webhooks, webhook)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel := connectingModel.(type) {
			case *WebhookWord:
				webhookWord := connectingModel
				webhookWord.ID = *webhook.ID
			case *WebhookUserMention:
				webhookUserMention := connectingModel
				webhookUserMention.WebhookID = *webhook.ID
			case *WebhookRoleMention:
				webhookRoleMention := connectingModel
				webhookRoleMention.WebhookID = *webhook.ID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhook)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook (
					webhook_serial_id,
					guild_id,
					webhook_id,
					subscription_type,
					subscription_id,
					last_posted_at
				) VALUES (
					:webhook_serial_id,
					:guild_id,
					:webhook_id,
					:subscription_type,
					:subscription_id,
					:last_posted_at
				)
			`, webhook)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
