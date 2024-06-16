package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestGetLineNgDiscordRoleID(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewLineNgDiscordRoleID(ctx, func(lng *fixtures.LineNgDiscordRoleID) {
			lng.ChannelID = "987654321"
			lng.GuildID = "123456789"
			lng.RoleID = "123456789"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからNG Discord IDを取得できること", func(t *testing.T) {
		ngDiscordIDs, err := repo.GetLineNgDiscordRoleIDByChannelID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", ngDiscordIDs[0])
	})
}

func TestRepository_InsertLineNgDiscordRoleIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく追加されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_role_id")

		repo := NewRepository(tx)
		lineNgDiscordIDs := []LineNgDiscordRoleIDAllCoulmns{
			{
				ChannelID: "123456789",
				GuildID:   "987654321",
				RoleID:    "123456789",
			},
			{
				ChannelID: "123456789",
				GuildID:   "123456789",
				RoleID:    "987654321",
			},
			{
				ChannelID: "987654321",
				GuildID:   "123456789",
				RoleID:    "987654321",
			},
		}
		err = repo.InsertLineNgDiscordRoleIDs(ctx, lineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_role_id")
		assert.NoError(t, err)

		assert.Equal(t, 3, lineChannelCount)
	})
}

func TestRepository_DeleteLineNgDiscordRoleIDs(t *testing.T) {
	ctx := context.Background()
	t.Run("NGなIDが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURLWithSslmode())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM line_ng_discord_role_id")

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewLineNgDiscordRoleID(ctx, func(lnt *fixtures.LineNgDiscordRoleID) {
				lnt.ChannelID = "123456789"
				lnt.GuildID = "987654321"
				lnt.RoleID = "123456789"
			}),
			fixtures.NewLineNgDiscordRoleID(ctx, func(lnt *fixtures.LineNgDiscordRoleID) {
				lnt.ChannelID = "123456789"
				lnt.GuildID = "123456789"
				lnt.RoleID = "987654321"
			}),
			fixtures.NewLineNgDiscordRoleID(ctx, func(lnt *fixtures.LineNgDiscordRoleID) {
				lnt.ChannelID = "987654321"
				lnt.GuildID = "123456789"
				lnt.RoleID = "123456789"
			}),
		)

		repo := NewRepository(tx)
		insertLineNgDiscordIDs := []LineNgDiscordRoleIDAllCoulmns{
			{
				ChannelID: "123456789",
				RoleID:    "123456789",
			},
			{
				ChannelID: "987654321",
				RoleID:    "123456789",
			},
		}

		err = repo.DeleteRoleIDsNotInProvidedList(ctx, "123456789", insertLineNgDiscordIDs)
		assert.NoError(t, err)

		var lineChannelCount int
		err = tx.GetContext(ctx, &lineChannelCount, "SELECT COUNT(*) FROM line_ng_discord_role_id")
		assert.NoError(t, err)

		assert.Equal(t, 2, lineChannelCount)
	})
}
