package fixtures

import (
	"context"
	"testing"
)

type WebhookUserMention struct {
	WebhookID int64  `db:"webhook_id"`
	UserID    string `db:"user_id"`
}

func NewWebhookUserMention(ctx context.Context, setter ...func(b *WebhookUserMention)) *ModelConnector {
	webhookUserMention := &WebhookUserMention{
		WebhookID: 1,
		UserID:    "1111111111111",
	}

	return &ModelConnector{
		Model: webhookUserMention,
		setter: func() {
			for _, s := range setter {
				s(webhookUserMention)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.WebhookUserMentions = append(f.WebhookUserMentions, webhookUserMention)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhookUserMention)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook_user_mention (
					webhook_id,
					user_id
				) VALUES (
					:webhook_id,
					:user_id
				)
			`, webhookUserMention)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
