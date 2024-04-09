package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestGetLinePostDiscordChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_post_discord_channel")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLinePostDiscordChannel(ctx, func(lc *fixtures.LinePostDiscordChannel) {
			lc.ChannelID = "123456789"
			lc.GuildID = "987654321"
			lc.Ng = false
			lc.BotMessage = false
		}),
	)
	repo := NewRepository(tx)
	t.Run("ChannelIDから送信しないかどうか取得できること", func(t *testing.T) {
		channel, err := repo.GetLinePostDiscordChannel(ctx, "123456789")
		assert.NoError(t, err)
		assert.Equal(t, false, channel.Ng)
		assert.Equal(t, false, channel.BotMessage)
	})
}

func TestInsertLinePostDiscordChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM line_post_discord_channel")

	repo := NewRepository(tx)

	var channels []LinePostDiscordChannelAllColumns
	t.Run("ChannelIDを追加できること", func(t *testing.T) {
		err := repo.InsertLinePostDiscordChannel(ctx, "123456789", "987654321")
		assert.NoError(t, err)
		query := `
			SELECT
				*
			FROM
				line_post_discord_channel
			WHERE
				channel_id = $1
		`
		err = tx.SelectContext(ctx, &channels, query, "123456789")
		assert.NoError(t, err)

		assert.Equal(t, "123456789", channels[0].ChannelID)
		assert.Equal(t, "987654321", channels[0].GuildID)
		assert.Equal(t, false, channels[0].Ng)
		assert.Equal(t, false, channels[0].BotMessage)
	})
}

func TestRepository_UpdateLinePostDiscordChannel(t *testing.T) {
	ctx := context.Background()
	t.Run("Channelが正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_post_discord_channel")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLinePostDiscordChannel(ctx, func(lc *fixtures.LinePostDiscordChannel) {
				lc.ChannelID = "123456789"
				lc.GuildID = "987654321"
				lc.Ng = false
				lc.BotMessage = false
			}),
		)

		repo := NewRepository(tx)
		updateLinePostDiscordChannel := LinePostDiscordChannelAllColumns{
			ChannelID:  "123456789",
			GuildID:    "987654321",
			Ng:         true,
			BotMessage: true,
		}
		err = repo.UpdateLinePostDiscordChannel(ctx, updateLinePostDiscordChannel)
		assert.NoError(t, err)

		var lineChannel LinePostDiscordChannelAllColumns
		err = tx.GetContext(ctx, &lineChannel, "SELECT * FROM line_post_discord_channel WHERE channel_id = $1", "123456789")
		assert.NoError(t, err)

		assert.Equal(t, "123456789", lineChannel.ChannelID)
		assert.Equal(t, "987654321", lineChannel.GuildID)
		assert.Equal(t, true, lineChannel.Ng)
		assert.Equal(t, true, lineChannel.BotMessage)
	})
}
