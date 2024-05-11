package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestGetLineNgDiscordMessageType(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgDiscordMessageType(ctx, func(lnt *fixtures.LineNgDiscordMessageType) {
			lnt.ChannelID = "987654321"
			lnt.Type = 6
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNGタイプを取得できること", func(t *testing.T) {
		ngTypes, err := repo.GetLineNgDiscordMessageTypeByChannelID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, []int{6}, ngTypes)
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

		_, err = tx.ExecContext(ctx, "DELETE FROM line_ng_discord_message_type")
		assert.NoError(t, err)

		repo := NewRepository(tx)
		lineNgDiscordTypes := []LineNgDiscordMessageType{
			{
				ChannelID: "123456789",
				GuildID:   "987654321",
				Type:      6,
			},
			{
				ChannelID: "123456789",
				GuildID:   "987654321",
				Type:      7,
			},
			{
				ChannelID: "987654321",
				GuildID:   "123456789",
				Type:      6,
			},
		}
		err = repo.InsertLineNgDiscordMessageTypes(ctx, lineNgDiscordTypes)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_message_type")
		assert.NoError(t, err)

		assert.Equal(t, 3, lineChannelCount)
	})
}

func TestRepository_DeleteLineNgDiscordMessageTypes(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなメッセージタイプが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		_, err = tx.ExecContext(ctx, "DELETE FROM line_ng_discord_message_type")
		assert.NoError(t, err)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineNgDiscordMessageType(ctx, func(lnt *fixtures.LineNgDiscordMessageType) {
				lnt.ChannelID = "123456789"
				lnt.Type = 6
			}),
			fixtures.NewLineNgDiscordMessageType(ctx, func(lnt *fixtures.LineNgDiscordMessageType) {
				lnt.ChannelID = "123456789"
				lnt.Type = 7
			}),
			fixtures.NewLineNgDiscordMessageType(ctx, func(lnt *fixtures.LineNgDiscordMessageType) {
				lnt.ChannelID = "987654321"
				lnt.Type = 6
			}),
		)

		repo := NewRepository(tx)
		insertLineNgDiscordTypes := []LineNgDiscordMessageType{
			{
				ChannelID: "123456789",
				Type:      6,
			},
			{
				ChannelID: "987654321",
				Type:      6,
			},
		}

		err = repo.DeleteMessageTypesNotInProvidedList(ctx, insertLineNgDiscordTypes)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_message_type")
		assert.NoError(t, err)

		assert.Equal(t, 2, lineChannelCount)
	})
}
