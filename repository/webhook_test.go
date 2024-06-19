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

func TestWebhook(t *testing.T) {
	ctx := context.Background()
	lastPostedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t.Run("Webhook登録", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		repo := NewRepository(tx)
		err = repo.InsertWebhook(
			ctx,
			"1111",
			"22222",
			"youtube",
			"test",
			lastPostedAt,
		)
		assert.NoError(t, err)

		var webhook Webhook
		err = tx.GetContext(ctx, &webhook, "SELECT * FROM webhook WHERE webhook_id = '22222'")
		assert.NoError(t, err)
		assert.Equal(t, "1111", webhook.GuildID)
		assert.Equal(t, "22222", webhook.WebhookID)
		assert.Equal(t, "youtube", webhook.SubscriptionType)
		assert.Equal(t, "test", webhook.SubscriptionID)
		assert.Equal(t, lastPostedAt, webhook.LastPostedAt.UTC())
	})

	t.Run("Webhook取得", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

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

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "1111")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 1)
		assert.Equal(t, "1111", webhooks[0].GuildID)
		assert.Equal(t, "22222", webhooks[0].WebhookID)
		assert.Equal(t, "youtube", webhooks[0].SubscriptionType)
		assert.Equal(t, "test", webhooks[0].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[0].LastPostedAt.UTC())
	})

	t.Run("Webhook取得_存在しない", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		repo := NewRepository(tx)

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "1111")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 0)
	})

	t.Run("Webhook取得_複数", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "33333"
				b.SubscriptionType = "niconico"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
		)

		repo := NewRepository(tx)

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "1111")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 2)
		assert.Equal(t, "1111", webhooks[0].GuildID)
		assert.Equal(t, "22222", webhooks[0].WebhookID)
		assert.Equal(t, "youtube", webhooks[0].SubscriptionType)
		assert.Equal(t, "test", webhooks[0].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[0].LastPostedAt.UTC())
		assert.Equal(t, "1111", webhooks[1].GuildID)
		assert.Equal(t, "33333", webhooks[1].WebhookID)
		assert.Equal(t, "niconico", webhooks[1].SubscriptionType)
		assert.Equal(t, "test", webhooks[1].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[1].LastPostedAt.UTC())
	})

	t.Run("Webhook取得_複数_存在しない", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "33333"
				b.SubscriptionType = "niconico"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
		)

		repo := NewRepository(tx)

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "2222")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 0)
	})

	t.Run("Webhook取得_複数_異なるguildID", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "33333"
				b.SubscriptionType = "niconico"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
		)

		repo := NewRepository(tx)

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "2222")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 0)
	})

	t.Run("Webhook取得_複数_異なるguildID_存在する", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "33333"
				b.SubscriptionType = "niconico"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}),
		)

		repo := NewRepository(tx)

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "1111")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 2)
		assert.Equal(t, "1111", webhooks[0].GuildID)
		assert.Equal(t, "22222", webhooks[0].WebhookID)
		assert.Equal(t, "youtube", webhooks[0].SubscriptionType)
		assert.Equal(t, "test", webhooks[0].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[0].LastPostedAt.UTC())
		assert.Equal(t, "1111", webhooks[1].GuildID)
		assert.Equal(t, "33333", webhooks[1].WebhookID)
		assert.Equal(t, "niconico", webhooks[1].SubscriptionType)
		assert.Equal(t, "test", webhooks[1].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[1].LastPostedAt.UTC())
	})

	t.Run("Webhook更新(最終更新日)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook")

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

		webhooks, err := repo.GetAllColumnsWebhooksByGuildID(ctx, "1111")
		assert.NoError(t, err)
		assert.Len(t, webhooks, 1)
		assert.Equal(t, "1111", webhooks[0].GuildID)
		assert.Equal(t, "22222", webhooks[0].WebhookID)
		assert.Equal(t, "youtube", webhooks[0].SubscriptionType)
		assert.Equal(t, "test", webhooks[0].SubscriptionID)
		assert.Equal(t, lastPostedAt, webhooks[0].LastPostedAt.UTC())

		err = repo.UpdateWebhookWithLastPostedAt(
			ctx,
			*f.Webhooks[0].WebhookSerialID,
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)

		var webhook Webhook
		err = tx.GetContext(ctx, &webhook, "SELECT * FROM webhook WHERE webhook_id = '22222'")
		assert.NoError(t, err)

		assert.Equal(t, "1111", webhook.GuildID)
		assert.Equal(t, "22222", webhook.WebhookID)
		assert.Equal(t, "youtube", webhook.SubscriptionType)
		assert.Equal(t, "test", webhook.SubscriptionID)
		assert.Equal(t, time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), webhook.LastPostedAt.UTC())
	})
}
