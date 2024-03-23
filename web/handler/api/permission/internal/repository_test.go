package internal

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/fixtures"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/web/config"

	"github.com/stretchr/testify/assert"
)

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
