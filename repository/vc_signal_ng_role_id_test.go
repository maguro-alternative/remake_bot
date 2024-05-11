package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
)

func TestInsertVcSignalNgRoleID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgRoleIDを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		repo := NewRepository(tx)
		err = repo.InsertVcSignalNgRole(ctx, "111","1111","11111")
		assert.NoError(t, err)

		var ngRole VcSignalNgRoleAllColumn
		err = tx.GetContext(ctx, &ngRole, "SELECT * FROM vc_signal_ng_role_id WHERE role_id = '11111'")
		assert.NoError(t, err)

		assert.Equal(t, "111", ngRole.VcChannelID)
		assert.Equal(t, "1111", ngRole.GuildID)
		assert.Equal(t, "11111", ngRole.RoleID)
	})

	t.Run("既にある場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_user_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.InsertVcSignalNgRole(ctx, "111","1111","11111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 1)
	})
}

func TestGetVcSignalNgRolesByChannelIDAllColumn(t *testing.T) {
	ctx := context.Background()
	t.Run("NgRoleIDを取得できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11112"
			}),
		)

		repo := NewRepository(tx)
		ngRoles, err := repo.GetVcSignalNgRolesByChannelIDAllColumn(ctx, "111")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 2)
	})

	t.Run("存在しない場合は空のスライスを返す", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		repo := NewRepository(tx)
		ngRoles, err := repo.GetVcSignalNgRolesByChannelIDAllColumn(ctx, "111")
		assert.NoError(t, err)

		assert.Len(t, ngRoles, 0)
	})
}

func TestDeleteVcSignalNgRoleID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgRoleIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})

	t.Run("存在しない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByChannelID(ctx, "111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})
}

func TestDeleteVcSignalNgRoleByGuildID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgRoleIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})

	t.Run("存在しない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByGuildID(ctx, "1111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})
}

func TestDeleteVcSignalNgRoleByRoleID(t *testing.T) {
	ctx := context.Background()
	t.Run("NgRoleIDを削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByRoleID(ctx, "11111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})

	t.Run("存在しない場合はエラーを返さずそのまま", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		repo := NewRepository(tx)
		err = repo.DeleteVcNgRoleByRoleID(ctx, "11111")
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 0)
	})
}

func TestDeleteRolesNotInProvidedList(t *testing.T) {
	ctx := context.Background()
	t.Run("指定されたNgRoleID以外を削除できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11112"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteRolesNotInProvidedList(ctx, "111", []string{"11111"})
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 1)
	})

	t.Run("指定されたNgRoleID以外を削除できること(すべて削除)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11112"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteRolesNotInProvidedList(ctx, "111", []string{"11111","11112"})
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 2)
	})

	t.Run("指定されたNgRoleID以外を削除できること(存在しないものがあっても無視)", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_ng_role_id")
		f := &fixtures.Fixture{DBv1: tx}
		f.Build(t,
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11111"
			}),
			fixtures.NewVcSignalNgRoleID(ctx, func(v *fixtures.VcSignalNgRoleID) {
				v.VcChannelID = "111"
				v.GuildID = "1111"
				v.RoleID = "11112"
			}),
		)

		repo := NewRepository(tx)
		err = repo.DeleteRolesNotInProvidedList(ctx, "111", []string{"11111","11112","11113"})
		assert.NoError(t, err)

		var ngRoles []VcSignalNgRoleAllColumn
		err = tx.SelectContext(ctx, &ngRoles, "SELECT * FROM vc_signal_ng_role_id")
		assert.NoError(t, err)
		assert.Len(t, ngRoles, 2)
	})


}
