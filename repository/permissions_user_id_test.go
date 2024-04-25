package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestRepository_InsertPermissionsID(t *testing.T) {
	ctx := context.Background()
	t.Run("Permissions_idが正しく挿入されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		repo := NewRepository(tx)
		insertPermissionsID := []PermissionUserIDAllColumns{
			{
				GuildID:    "987654321",
				Type:       "line_bot",
				UserID:     "123456789",
				Permission: "all",
			},
		}
		err = repo.InsertPermissionUserIDs(ctx, insertPermissionsID)
		assert.NoError(t, err)

		var permissionsID PermissionUserIDAllColumns
		err = tx.GetContext(ctx, &permissionsID, "SELECT * FROM permissions_user_id WHERE guild_id = $1", "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", permissionsID.GuildID)
		assert.Equal(t, "line_bot", permissionsID.Type)
		assert.Equal(t, "123456789", permissionsID.UserID)
		assert.Equal(t, "all", permissionsID.Permission)
	})
}

func TestGetGuildPermissionIDsAllColumns(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "123456789"
			p.Type = "line_bot"
			p.Permission = "read"
		}),
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "345678912"
			p.Type = "line_bot"
			p.Permission = "write"
		}),
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "567891234"
			p.Type = "line_bot"
			p.Permission = "all"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionIDを取得できること", func(t *testing.T) {
		permissionIDs, err := repo.GetGuildPermissionUserIDsAllColumns(ctx, "987654321")
		assert.NoError(t, err)
		assert.Equal(t, "987654321", permissionIDs[0].GuildID)
		assert.Equal(t, "line_bot", permissionIDs[0].Type)
		assert.Equal(t, "123456789", permissionIDs[0].UserID)
		assert.Equal(t, "read", permissionIDs[0].Permission)
		assert.Equal(t, "987654321", permissionIDs[1].GuildID)
		assert.Equal(t, "line_bot", permissionIDs[1].Type)
		assert.Equal(t, "345678912", permissionIDs[1].UserID)
		assert.Equal(t, "write", permissionIDs[1].Permission)
		assert.Equal(t, "987654321", permissionIDs[2].GuildID)
		assert.Equal(t, "line_bot", permissionIDs[2].Type)
		assert.Equal(t, "567891234", permissionIDs[2].UserID)
		assert.Equal(t, "all", permissionIDs[2].Permission)
	})
}

func TestGetPermissionIDs(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	assert.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	assert.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "123456789"
			p.Type = "line_bot"
			p.Permission = "read"
		}),
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "345678912"
			p.Type = "line_bot"
			p.Permission = "write"
		}),
		fixtures.NewPermissionsUserID(ctx, func(p *fixtures.PermissionsUserID) {
			p.GuildID = "987654321"
			p.UserID = "567891234"
			p.Type = "line_bot"
			p.Permission = "all"
		}),
	)
	repo := NewRepository(tx)
	t.Run("GuildIDからPermissionIDを取得できること", func(t *testing.T) {
		permissionIDs, err := repo.GetPermissionUserIDs(ctx, "987654321", "line_bot")
		assert.NoError(t, err)
		assert.Equal(t, "123456789", permissionIDs[0].UserID)
		assert.Equal(t, "read", permissionIDs[0].Permission)
		assert.Equal(t, "345678912", permissionIDs[1].UserID)
		assert.Equal(t, "write", permissionIDs[1].Permission)
		assert.Equal(t, "567891234", permissionIDs[2].UserID)
		assert.Equal(t, "all", permissionIDs[2].Permission)
	})
}

func TestRepository_DeletePermissionsID(t *testing.T) {
	ctx := context.Background()
	t.Run("Permissions_idが正しく削除されること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewPermissionsUserID(ctx, func(pi *fixtures.PermissionsUserID) {
				pi.GuildID = "987654321"
				pi.Type = "line_bot"
				pi.UserID = "123456789"
				pi.Permission = "all"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeletePermissionUserIDs(ctx, "987654321")
		assert.NoError(t, err)

		var permissionsID PermissionUserID
		err = tx.GetContext(ctx, &permissionsID, "SELECT * FROM permissions_user_id WHERE guild_id = $1", "987654321")
		assert.Error(t, err)
		assert.Empty(t, permissionsID)
	})
}
