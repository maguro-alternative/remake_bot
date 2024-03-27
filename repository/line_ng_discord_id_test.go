package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"

	"github.com/stretchr/testify/assert"
)

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


func TestRepository_InsertLineNgDiscordIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく追加されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_id")

		repo := NewRepository(tx)
		lineNgDiscordIDs := []LineNgDiscordIDAllCoulmns{
			{
				ChannelID: "123456789",
				GuildID:   "987654321",
				ID:        "123456789",
				IDType:    "user",
			},
			{
				ChannelID: "123456789",
				GuildID:   "123456789",
				ID:        "987654321",
				IDType:    "user",
			},
			{
				ChannelID: "987654321",
				GuildID:   "123456789",
				ID:        "987654321",
				IDType:    "user",
			},
		}
		err = repo.InsertLineNgDiscordIDs(ctx, lineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_id")
		assert.NoError(t, err)

		assert.Equal(t, 3, lineChannelCount)
	})
}

func TestRepository_DeleteLineNgDiscordIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_id")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineNgDiscordID(ctx, func(lnt *fixtures.LineNgDiscordID) {
				lnt.ChannelID = "123456789"
				lnt.ID = "123456789"
				lnt.IDType = "user"
			}),
			fixtures.NewLineNgDiscordID(ctx, func(lnt *fixtures.LineNgDiscordID) {
				lnt.ChannelID = "123456789"
				lnt.ID = "987654321"
				lnt.IDType = "user"
			}),
			fixtures.NewLineNgDiscordID(ctx, func(lnt *fixtures.LineNgDiscordID) {
				lnt.ChannelID = "987654321"
				lnt.ID = "123456789"
				lnt.IDType = "user"
			}),
		)

		repo := NewRepository(tx)
		insertLineNgDiscordIDs := []LineNgDiscordIDAllCoulmns{
			{
				ChannelID: "123456789",
				ID:        "123456789",
				IDType:    "user",
			},
			{
				ChannelID: "987654321",
				ID:        "123456789",
				IDType:    "user",
			},
		}

		err = repo.DeleteNotInsertLineNgDiscordIDs(ctx, insertLineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_id")
		assert.NoError(t, err)

		assert.Equal(t, 2, lineChannelCount)
	})
}

