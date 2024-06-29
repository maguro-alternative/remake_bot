package fixtures

import (
	"context"
	"testing"
)

type WebhookWord struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	Condition       string `db:"conditions"`
	Word            string `db:"word"`
}

func NewWebhookWord(ctx context.Context, setter ...func(b *WebhookWord)) *ModelConnector {
	webhookWord := &WebhookWord{
		WebhookSerialID: 1,
		Condition:       "NgOr",
		Word:            "word",
	}

	return &ModelConnector{
		Model: webhookWord,
		setter: func() {
			for _, s := range setter {
				s(webhookWord)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.WebhookWords = append(f.WebhookWords, webhookWord)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhookWord)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook_word (
					webhook_serial_id,
					conditions,
					word
				) VALUES (
					:webhook_serial_id,
					:conditions,
					:word
				)
			`, webhookWord)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
