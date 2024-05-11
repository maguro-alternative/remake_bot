package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestGetLineNgDiscordUserID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgDiscordUserID(ctx, func(lng *fixtures.LineNgDiscordUserID) {
			lng.ChannelID = "987654321"
			lng.GuildID = "123456789"
			lng.UserID = "123456789"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNG Discord User IDを取得できること", func(t *testing.T) {
		ngDiscordIDs, err := repo.GetLineNgDiscordUserIDByChannelID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", ngDiscordIDs[0])
	})
}

func TestRepository_InsertLineNgDiscordUserIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく追加されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_user_id")

		repo := NewRepository(tx)
		lineNgDiscordIDs := []LineNgDiscordUserIDAllCoulmns{
			{
				ChannelID: "123456789",
				GuildID:   "987654321",
				UserID:    "123456789",
			},
			{
				ChannelID: "123456789",
				GuildID:   "123456789",
				UserID:    "987654321",
			},
			{
				ChannelID: "987654321",
				GuildID:   "123456789",
				UserID:    "987654321",
			},
		}
		err = repo.InsertLineNgDiscordUserIDs(ctx, lineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_user_id")
		assert.NoError(t, err)

		assert.Equal(t, 3, lineChannelCount)
	})
}

func TestRepository_DeleteLineNgDiscordUserIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_user_id")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineNgDiscordUserID(ctx, func(lnt *fixtures.LineNgDiscordUserID) {
				lnt.ChannelID = "123456789"
				lnt.UserID = "123456789"
			}),
			fixtures.NewLineNgDiscordUserID(ctx, func(lnt *fixtures.LineNgDiscordUserID) {
				lnt.ChannelID = "123456789"
				lnt.UserID = "987654321"
			}),
			fixtures.NewLineNgDiscordUserID(ctx, func(lnt *fixtures.LineNgDiscordUserID) {
				lnt.ChannelID = "987654321"
				lnt.UserID = "123456789"
			}),
		)

		repo := NewRepository(tx)
		insertLineNgDiscordIDs := []LineNgDiscordUserIDAllCoulmns{
			{
				ChannelID: "123456789",
				UserID:    "123456789",
			},
			{
				ChannelID: "987654321",
				UserID:    "123456789",
			},
		}

		err = repo.DeleteUserIDsNotInProvidedList(ctx, insertLineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_user_id")
		assert.NoError(t, err)

		assert.Equal(t, 2, lineChannelCount)
	})
}
