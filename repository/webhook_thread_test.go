package repository

import (
	"context"
	"testing"
	"time"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestWebhookThread(t *testing.T) {
	ctx := context.Background()
	lastPostedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t.Run("WebhookThread挿入", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_thread")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
		)

		repo := NewRepository(tx)

		err = repo.InsertWebhookThread(ctx, *f.Webhooks[0].WebhookSerialID, "test")
		assert.NoError(t, err)

		var webhookThreads []*WebhookThread
		err = tx.SelectContext(ctx, &webhookThreads, "SELECT * FROM webhook_thread")
		assert.NoError(t, err)
		assert.Len(t, webhookThreads, 1)
	})

	t.Run("WebhookThread挿入(既存のものはスルー)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_thread")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookThread(ctx, func(b *fixtures.WebhookThread) {
					b.WebhookSerialID = *f.Webhooks[0].WebhookSerialID
					b.ThreadID = "test"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.InsertWebhookThread(ctx, *f.Webhooks[0].WebhookSerialID, "test")
		assert.NoError(t, err)

		var webhookThreads []*WebhookThread
		err = tx.SelectContext(ctx, &webhookThreads, "SELECT * FROM webhook_thread")
		assert.NoError(t, err)
		assert.Len(t, webhookThreads, 1)
	})

	t.Run("WebhookThread削除", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_thread")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookThread(ctx, func(b *fixtures.WebhookThread) {
					b.WebhookSerialID = *f.Webhooks[0].WebhookSerialID
					b.ThreadID = "test"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.DeleteWebhookThreadsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID)
		assert.NoError(t, err)

		var webhookThreads []*WebhookThread
		err = tx.SelectContext(ctx, &webhookThreads, "SELECT * FROM webhook_thread")
		assert.NoError(t, err)
		assert.Len(t, webhookThreads, 0)
	})
}
