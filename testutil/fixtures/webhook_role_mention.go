package fixtures

import (
	"context"
	"testing"
)

type WebhookRoleMention struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	RoleID          string `db:"role_id"`
}

func NewWebhookMention(ctx context.Context, setter ...func(b *WebhookRoleMention)) *ModelConnector {
	webhookRoleMention := &WebhookRoleMention{
		WebhookSerialID: 1,
		RoleID:          "1111111111111",
	}

	return &ModelConnector{
		Model: webhookRoleMention,
		setter: func() {
			for _, s := range setter {
				s(webhookRoleMention)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.WebhookRoleMentions = append(f.WebhookRoleMentions, webhookRoleMention)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhookRoleMention)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook_role_mention (
					webhook_serial_id,
					role_id
				) VALUES (
					:webhook_serial_id,
					:role_id
				)
			`, webhookRoleMention)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
