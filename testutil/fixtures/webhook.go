package fixtures

import (
	"context"
	"testing"
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
				webhookWord.WebhookSerialID = *webhook.WebhookSerialID
			case *WebhookUserMention:
				webhookUserMention := connectingModel
				webhookUserMention.WebhookSerialID = *webhook.WebhookSerialID
			case *WebhookRoleMention:
				webhookRoleMention := connectingModel
				webhookRoleMention.WebhookSerialID = *webhook.WebhookSerialID
			case *WebhookThread:
				webhookThread := connectingModel
				webhookThread.WebhookSerialID = *webhook.WebhookSerialID
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhook)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			err := f.DBv1.QueryRowxContext(ctx, `
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
				) RETURNING webhook_serial_id`,
				webhook.GuildID,
				webhook.WebhookID,
				webhook.SubscriptionType,
				webhook.SubscriptionID,
				webhook.LastPostedAt,
			).Scan(&webhook.WebhookSerialID)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
