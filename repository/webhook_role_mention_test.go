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

func TestWebhookRoleMention(t *testing.T) {
	ctx := context.Background()
	lastPostedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t.Run("WebhookRoleMention挿入", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_role_mention")

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

		err = repo.InsertWebhookRoleMention(ctx, *f.Webhooks[0].WebhookSerialID, "111111")
		assert.NoError(t, err)

		var webhookRoleMentions []*WebhookRoleMention
		err = tx.SelectContext(ctx, &webhookRoleMentions, "SELECT * FROM webhook_role_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookRoleMentions, 1)
	})

	t.Run("WebhookRoleMention挿入(既存のものはスルー)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_role_mention")

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

		err = repo.InsertWebhookRoleMention(ctx, *f.Webhooks[0].WebhookSerialID, "111111")
		assert.NoError(t, err)

		err = repo.InsertWebhookRoleMention(ctx, *f.Webhooks[0].WebhookSerialID, "111111")
		assert.NoError(t, err)

		var webhookRoleMentions []*WebhookRoleMention
		err = tx.SelectContext(ctx, &webhookRoleMentions, "SELECT * FROM webhook_role_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookRoleMentions, 1)
	})

	t.Run("WebhookRoleMention取得(idsから取得)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_role_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookRoleMention(ctx, func(b *fixtures.WebhookRoleMention) {
					b.RoleID = "111111"
				}),
			),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22223"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookRoleMention(ctx, func(b *fixtures.WebhookRoleMention) {
					b.RoleID = "111111"
				}),
			),
		)

		repo := NewRepository(tx)

		webhookRoleMentions, err := repo.GetWebhookRoleMentionWithWebhookSerialIDs(ctx, []int64{*f.Webhooks[0].WebhookSerialID, *f.Webhooks[1].WebhookSerialID})
		assert.NoError(t, err)
		assert.Len(t, webhookRoleMentions, 2)
		assert.Equal(t, "111111", webhookRoleMentions[0].RoleID)
		assert.Equal(t, "111111", webhookRoleMentions[1].RoleID)
	})


	t.Run("WebhookRoleMention削除", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_role_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookRoleMention(ctx, func(b *fixtures.WebhookRoleMention) {
					b.RoleID = "111111"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.DeleteWebhookRoleMentionsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, []string{})
		assert.NoError(t, err)

		var webhookRoleMentions []*WebhookRoleMention
		err = tx.SelectContext(ctx, &webhookRoleMentions, "SELECT * FROM webhook_role_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookRoleMentions, 0)
	})

	t.Run("WebhookRoleMention削除(存在しないものはスルー)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_role_mention")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookRoleMention(ctx, func(b *fixtures.WebhookRoleMention) {
					b.RoleID = "111111"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.DeleteWebhookRoleMentionsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, []string{"111111", "222222"})
		assert.NoError(t, err)

		var webhookRoleMentions []*WebhookRoleMention
		err = tx.SelectContext(ctx, &webhookRoleMentions, "SELECT * FROM webhook_role_mention")
		assert.NoError(t, err)
		assert.Len(t, webhookRoleMentions, 1)
	})
}
