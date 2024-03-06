package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/stretchr/testify/assert"
)

func TestGetLineChannel(t *testing.T) {
	ctx := context.Background()
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
	t.Run("ChannelIDから送信しないかどうか取得できること", func(t *testing.T) {
		channel, err := repo.GetLineChannel(ctx, "123456789")
		assert.NoError(t, err)
		assert.Equal(t, false, channel.Ng)
		assert.Equal(t, false, channel.BotMessage)
	})
}

func TestInsertLineChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	repo := NewRepository(tx)

	var channels []TestLineChannel
	t.Run("ChannelIDを追加できること", func(t *testing.T) {
		err := repo.InsertLineChannel(ctx, "123456789", "987654321")
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

func TestGetLineNgType(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgType(ctx, func(lnt *fixtures.LineNgType) {
			lnt.ChannelID = "987654321"
			lnt.Type = 6
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNGタイプを取得できること", func(t *testing.T) {
		ngTypes, err := repo.GetLineNgType(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []int{6}, ngTypes)
	})
}

func TestGetLineNgDiscordID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgDiscordID(ctx, func(lng *fixtures.LineNgDiscordID) {
			lng.ChannelID = "987654321"
			lng.GuildID = "123456789"
			lng.IDType = "user"
			lng.ID = "123456789"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNG Discord IDを取得できること", func(t *testing.T) {
		ngDiscordIDs, err := repo.GetLineNgDiscordID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", ngDiscordIDs[0].ID)
		assert.Equal(t, "user", ngDiscordIDs[0].IDType)
	})
}
