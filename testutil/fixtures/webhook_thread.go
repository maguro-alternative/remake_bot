package fixtures

import (
	"context"
	"testing"
)

type WebhookThread struct {
	WebhookSerialID int64  `db:"webhook_serial_id"`
	ThreadID        string `db:"thread_id"`
}

func NewWebhookThread(ctx context.Context, setter ...func(b *WebhookThread)) *ModelConnector {
	webhookThread := &WebhookThread{
		WebhookSerialID: 1,
		ThreadID:        "thread_id",
	}

	return &ModelConnector{
		Model: webhookThread,
		setter: func() {
			for _, s := range setter {
				s(webhookThread)
			}
		},
		addToFixture: func(t *testing.T, f *Fixture) {
			f.WebhookThreads = append(f.WebhookThreads, webhookThread)
		},
		connect: func(t *testing.T, f *Fixture, connectingModel interface{}) {
			switch connectingModel.(type) {
			default:
				t.Fatalf("%T cannot be connected to %T", connectingModel, webhookThread)
			}
		},
		insertTable: func(t *testing.T, f *Fixture) {
			_, err := f.DBv1.NamedExecContext(ctx, `
				INSERT INTO webhook_thread (
					webhook_serial_id,
					thread_id
				) VALUES (
					:webhook_serial_id,
					:thread_id
				)
			`, webhookThread)
			if err != nil {
				t.Fatalf("insert error: %v", err)
			}
		},
	}
}
