package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/stretchr/testify/assert"
)

func TestRepository_UpdateLineBot(t *testing.T) {
	ctx := context.Background()
	t.Run("Channelが正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineChannel(ctx, func(lc *fixtures.LineChannel) {
				lc.ChannelID = "123456789"
				lc.GuildID = "987654321"
				lc.Ng = false
				lc.BotMessage = false
			}),
		)

		repo := NewRepository(tx)
		updateLineChannel := LineChannel{
				ChannelID:  "123456789",
				GuildID:    "987654321",
				Ng:         true,
				BotMessage: true,
		}
		err = repo.UpdateLinePostDiscordChannel(ctx, updateLineChannel)
		assert.NoError(t, err)

		var lineChannel LineChannel
		err = tx.GetContext(ctx, &lineChannel, "SELECT * FROM line_post_discord_channel WHERE channel_id = $1", "123456789")
		assert.NoError(t, err)

		assert.Equal(t, "123456789", lineChannel.ChannelID)
		assert.Equal(t, "987654321", lineChannel.GuildID)
		assert.Equal(t, true, lineChannel.Ng)
		assert.Equal(t, true, lineChannel.BotMessage)
	})
}

func TestRepository_InsertLineNgDiscordMessageTypes(t *testing.T) {
	ctx := context.Background()
	t.Run("Channelが正しく追加されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		repo := NewRepository(tx)
		lineNgDiscordTypes := []LineNgType{
			{
				ChannelID: "123456789",
				Type:      6,
			},
			{
				ChannelID: "123456789",
				Type:      7,
			},
			{
				ChannelID: "987654321",
				Type:      6,
			},
		}
		err = repo.InsertLineNgDiscordMessageTypes(ctx, lineNgDiscordTypes)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_post_discord_channel WHERE channel_id = $1", "123456789")
		assert.NoError(t, err)

		assert.Equal(t, 3, lineChannelCount)
	})
}

func TestRepository_DeleteLineNgDiscordMessageTypes(t *testing.T) {
	ctx := context.Background()
	t.Run("Channelが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineNgType(ctx, func(lnt *fixtures.LineNgType) {
				lnt.ChannelID = "123456789"
				lnt.Type = 6
			}),
			fixtures.NewLineNgType(ctx, func(lnt *fixtures.LineNgType) {
				lnt.ChannelID = "123456789"
				lnt.Type = 7
			}),
			fixtures.NewLineNgType(ctx, func(lnt *fixtures.LineNgType) {
				lnt.ChannelID = "987654321"
				lnt.Type = 6
			}),
		)

		repo := NewRepository(tx)
		insertLineNgDiscordTypes := []LineNgType{
			{
				ChannelID: "123456789",
				Type:      6,
			},
			{
				ChannelID: "987654321",
				Type:      6,
			},
		}

		err = repo.DeleteNotInsertLineNgDiscordMessageTypes(ctx, insertLineNgDiscordTypes)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_post_discord_channel")
		assert.NoError(t, err)

		assert.Equal(t, 2, lineChannelCount)
	})
}
