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

func TestWebhookWord(t *testing.T) {
	ctx := context.Background()
	lastPostedAt := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	t.Run("WebhookWord挿入", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

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

		err = repo.InsertWebhookWord(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", "word")
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 1)
	})

	t.Run("WebhookWord挿入(既存のものはスルー)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
			),
		)

		repo := NewRepository(tx)

		err = repo.InsertWebhookWord(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", "word")
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 1)
	})

	t.Run("WebhookWord取得", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
			),
		)

		webhookWords, err := repo.GetWebhookWordWithWebhookSerialIDAndCondition(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 1)
		assert.Equal(t, "ng_or", webhookWords[0].Condition)
		assert.Equal(t, "word", webhookWords[0].Word)
	})

	t.Run("WebhookWord取得(指定されたidのみ取得)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "search_or"
					b.Word = "word3"
				}),
			),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22223"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test2"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "search_or"
					b.Word = "word3"
				}),
			),
		)

		webhookWords, err := repo.GetWebhookWordWithWebhookSerialIDs(ctx, []int64{*f.Webhooks[0].WebhookSerialID, *f.Webhooks[1].WebhookSerialID})
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 6)
	})

	t.Run("WebhookWord削除", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word2"
				}),
			),
		)

		err = repo.DeleteWebhookWordsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", []string{"word"})
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 1)
	})

	t.Run("WebhookWord削除(ng_or以外は削除しない)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
			),
		)

		err = repo.DeleteWebhookWordsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", []string{"word"})
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 2)
	})

	t.Run("WebhookWord削除(削除はしない)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
			),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22223"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
			),
		)

		err = repo.DeleteWebhookWordsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", []string{"word"})
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 2)
	})

	t.Run("WebhookWord削除(はじめのword2を削除)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()

		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM webhook_word")

		repo := NewRepository(tx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22222"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word2"
				}),
			),
			fixtures.NewWebhook(ctx, func(b *fixtures.Webhook) {
				b.GuildID = "1111"
				b.WebhookID = "22223"
				b.SubscriptionType = "youtube"
				b.SubscriptionID = "test"
				b.LastPostedAt = lastPostedAt
			}).Connect(
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_or"
					b.Word = "word"
				}),
				fixtures.NewWebhookWord(ctx, func(b *fixtures.WebhookWord) {
					b.Condition = "ng_and"
					b.Word = "word2"
				}),
			),
		)

		err = repo.DeleteWebhookWordsNotInProvidedList(ctx, *f.Webhooks[0].WebhookSerialID, "ng_or", []string{"word"})
		assert.NoError(t, err)

		var webhookWords []*WebhookWord
		err = tx.SelectContext(ctx, &webhookWords, "SELECT * FROM webhook_word")
		assert.NoError(t, err)
		assert.Len(t, webhookWords, 3)
	})
}
