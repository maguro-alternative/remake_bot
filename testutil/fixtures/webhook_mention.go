package fixtures

import (
	"context"
	"testing"
)

type WebhookMention struct {
	WebhookID int64  `db:"webhook_id"`
	IDType    string `db:"id_type"`
	ID        string `db:"id"`
}

func NewWebhookMention(ctx context.Context, setter ...func(b *WebhookMention)) *ModelConnector {
	webhookMention := &WebhookMention{
		WebhookID: 1,
		IDType:    "user",
		ID:        "1111111111111",
	}

	return &ModelConnector{
		Model: webhookMention,
		setter: func() {
			for _, s := range setter {
				s(webhookMention)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.WebhookMentions = append(f.WebhookMentions, webhookMention)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhookMention)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook_mention (
					webhook_id,
					id_type,
					id
				) VALUES (
					:webhook_id,
					:id_type,
					:id
				)
			`, webhookMention)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
