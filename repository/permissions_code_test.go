package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestGetPermissionsCode(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsCode(ctx, func(p *fixtures.PermissionsCode) {
			p.GuildID = "987654321"
			p.Type = "line_bot"
			p.Code = 8
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionsCodeを取得できること", func(t *testing.T) {
		permissionCode, err := repo.GetPermissionCodeByGuildIDAndType(ctx, "987654321", "line_bot")
		assert.NoError(t, err)
		assert.Equal(t, int64(8), permissionCode)
	})
}

func TestGetPermissionsCodes(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsCode(ctx, func(p *fixtures.PermissionsCode) {
			p.GuildID = "987654321"
			p.Type = "line_bot"
			p.Code = 8
		}),
		fixtures.NewPermissionsCode(ctx, func(p *fixtures.PermissionsCode) {
			p.GuildID = "987654321"
			p.Type = "line_post_discord_channel"
			p.Code = 9
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionsCodeを取得できること", func(t *testing.T) {
		permissionCode, err := repo.GetPermissionCodesByGuildID(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, 2, len(permissionCode))
		assert.Equal(t, int64(8), permissionCode[0].Code)
		assert.Equal(t, int64(9), permissionCode[1].Code)
	})
}

func TestRepository_UpdatePermissionsCode(t *testing.T) {
	ctx := context.Background()
	t.Run("Permissions_codeが正しく更新されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewPermissionsCode(ctx, func(pc *fixtures.PermissionsCode) {
				pc.GuildID = "987654321"
				pc.Type = "line_bot"
				pc.Code = 8
			}),
		)

		repo := NewRepository(tx)
		updatePermissionsCode := []PermissionCode{
			{
				GuildID: "987654321",
				Type:    "line_bot",
				Code:    9,
			},
		}
		err = repo.UpdatePermissionCodes(ctx, updatePermissionsCode)
		assert.NoError(t, err)

		var permissionsCode PermissionCode
		err = tx.GetContext(ctx, &permissionsCode, "SELECT * FROM permissions_code WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", permissionsCode.GuildID)
		assert.Equal(t, "line_bot", permissionsCode.Type)
		assert.Equal(t, int64(9), permissionsCode.Code)
	})
}
