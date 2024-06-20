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

func TestWebhookUserMention(t *testing.T) {
	ctx := context.Background()
	lastPostedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t.Run("WebhookUserMention挿入", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_user_mention")

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

		err = repo.InsertWebhookUserMention(ctx, *f.Webhooks[0].WebhookSerialID, "111111")
		assert.NoError(t, err)

		var webhookUserMentions []*WebhookUserMention
		err = tx.SelectContext(ctx, &webhookUserMentions, "SELECT * FROM webhook_user_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookUserMentions, 1)
	})

	t.Run("WebhookUserMention挿入(既存のものはスルー)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_user_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookUserMention(ctx, func(b *fixtures.WebhookUserMention) {
					b.UserID = "111111"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.InsertWebhookUserMention(ctx, *f.Webhooks[0].WebhookSerialID, "111111")
		assert.NoError(t, err)

		var webhookUserMentions []*WebhookUserMention
		err = tx.SelectContext(ctx, &webhookUserMentions, "SELECT * FROM webhook_user_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookUserMentions, 1)
	})

	t.Run("WebhookUserMention取得", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_user_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookUserMention(ctx, func(b *fixtures.WebhookUserMention) {
					b.UserID = "111111"
				}),
			),
		)

		repo := NewRepository(tx)

		webhookUserMentions, err := repo.GetWebhookUserMentionWithWebhookSerialID(ctx, *f.Webhooks[0].WebhookSerialID)
		assert.NoError(t, err)
		assert.Len(t, webhookUserMentions, 1)
		assert.Equal(t, "111111", webhookUserMentions[0].UserID)
	})

	t.Run("WebhookUserMention削除(指定したもの以外は削除される)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_user_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookUserMention(ctx, func(b *fixtures.WebhookUserMention) {
					b.UserID = "111111"
				}),
				fixtures.NewWebhookUserMention(ctx, func(b *fixtures.WebhookUserMention) {
					b.UserID = "222222"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.DeleteWebhookUserMentionsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, []string{"111111"})
		assert.NoError(t, err)

		var webhookUserMentions []*WebhookUserMention
		err = tx.SelectContext(ctx, &webhookUserMentions, "SELECT * FROM webhook_user_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookUserMentions, 1)
	})
}
